package sso

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"backend-alusi-go/config"
)

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	IDToken     string `json:"id_token,omitempty"`
	ExpiresIn   int    `json:"expires_in,omitempty"`
}

type UserInfo struct {
	Sub           string `json:"sub"`
	UserType      string `json:"user_type"` // internal (Pegawai) or external (Mitra)
	NamaLengkap   string `json:"nama_lengkap"`
	Email         string `json:"email"`
	NIP           string `json:"nip,omitempty"`
	NIPLama       string `json:"nip_lama,omitempty"`
	KodeSatker    string `json:"kode_satker,omitempty"`
	NIK           string `json:"nik,omitempty"`
	KodeKabupaten string `json:"kode_kabupaten,omitempty"`
	KodeKecamatan string `json:"kode_kecamatan,omitempty"`
	KodeDesa      string `json:"kode_desa,omitempty"`
	KodeSLS       string `json:"kode_sls,omitempty"`
	Jabatan       string `json:"jabatan,omitempty"`
}

type Client struct {
	cfg        *config.SSOConfig
	httpClient *http.Client
}

func NewClient(cfg *config.SSOConfig) *Client {
	return &Client{
		cfg: cfg,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetAuthURL generates authorization URL for initiating SSO login
func (c *Client) GetAuthURL(state string) string {
	params := url.Values{}
	params.Add("client_id", c.cfg.ClientID)
	params.Add("redirect_uri", c.cfg.RedirectURI)
	params.Add("response_type", "code")
	params.Add("scope", "openid profile email")
	params.Add("state", state)

	baseURL := strings.TrimRight(c.cfg.IssuerURL, "/")
	return fmt.Sprintf("%s/authorize?%s", baseURL, params.Encode())
}

// ExchangeCode exchanges authorization code for access token via POST /token
func (c *Client) ExchangeCode(ctx context.Context, code string) (*TokenResponse, error) {
	form := url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"client_id":     {c.cfg.ClientID},
		"client_secret": {c.cfg.ClientSecret},
		"redirect_uri":  {c.cfg.RedirectURI},
	}

	tokenURL := fmt.Sprintf("%s/token", strings.TrimRight(c.cfg.IssuerURL, "/"))
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create token request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("token request failed: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read token response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token endpoint returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var tok TokenResponse
	if err := json.Unmarshal(bodyBytes, &tok); err != nil {
		return nil, fmt.Errorf("failed to decode token response: %w", err)
	}

	if tok.AccessToken == "" {
		return nil, fmt.Errorf("empty access_token received from SSO")
	}

	return &tok, nil
}

// GetUserInfo retrieves user profile claims via GET /userinfo or JWT token claims
func (c *Client) GetUserInfo(ctx context.Context, accessToken string, idToken ...string) (*UserInfo, error) {
	baseURL := strings.TrimRight(c.cfg.IssuerURL, "/")

	// Possible userinfo endpoints on different OIDC providers
	endpoints := []string{
		baseURL + "/userinfo",
		baseURL + "/api/userinfo",
		baseURL + "/api/v1/userinfo",
		baseURL + "/oauth/userinfo",
		baseURL + "/me",
	}

	var lastErr error
	for _, userInfoURL := range endpoints {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, userInfoURL, nil)
		if err != nil {
			continue
		}
		req.Header.Set("Authorization", "Bearer "+accessToken)

		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = err
			continue
		}
		bodyBytes, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			var user UserInfo
			if err := json.Unmarshal(bodyBytes, &user); err == nil && user.Sub != "" {
				return &user, nil
			}
		} else {
			lastErr = fmt.Errorf("endpoint %s returned status %d", userInfoURL, resp.StatusCode)
		}
	}

	// Fallback 1: Parse claims from IDToken if provided
	if len(idToken) > 0 && idToken[0] != "" {
		if claims, err := parseJWTClaims(idToken[0]); err == nil && claims.Sub != "" {
			return claims, nil
		}
	}

	// Fallback 2: Parse claims directly from AccessToken if it's a JWT
	if claims, err := parseJWTClaims(accessToken); err == nil && claims.Sub != "" {
		return claims, nil
	}

	if lastErr != nil {
		return nil, fmt.Errorf("failed to fetch userinfo from SSO provider: %w", lastErr)
	}

	return nil, fmt.Errorf("endpoint userinfo SSO (404) dan klaim token tidak ditemukan")
}

// parseJWTClaims extracts user profile claims from a JWT token string (id_token or access_token)
func parseJWTClaims(tokenStr string) (*UserInfo, error) {
	if tokenStr == "" {
		return nil, fmt.Errorf("token string is empty")
	}

	parts := strings.Split(tokenStr, ".")
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid JWT format")
	}

	segment := parts[1]
	switch len(segment) % 4 {
	case 2:
		segment += "=="
	case 3:
		segment += "="
	}

	decoded, err := base64.URLEncoding.DecodeString(segment)
	if err != nil {
		decoded, err = base64.RawURLEncoding.DecodeString(parts[1])
		if err != nil {
			return nil, fmt.Errorf("failed to base64 decode JWT payload: %w", err)
		}
	}

	var claims struct {
		Sub           string `json:"sub"`
		UserType      string `json:"user_type"`
		Nama          string `json:"nama"`
		NamaLengkap   string `json:"nama_lengkap"`
		Name          string `json:"name"`
		Username      string `json:"username"`
		Email         string `json:"email"`
		NIP           string `json:"nip"`
		NIPLama       string `json:"nip_lama"`
		NipLama       string `json:"niplama"`
		KodeSatker    string `json:"kode_satker"`
		NIK           string `json:"nik"`
		KodeKabupaten string `json:"kode_kabupaten"`
		KodeKecamatan string `json:"kode_kecamatan"`
		KodeDesa      string `json:"kode_desa"`
		KodeSLS       string `json:"kode_sls"`
		Jabatan       string `json:"jabatan"`
	}

	if err := json.Unmarshal(decoded, &claims); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JWT payload: %w", err)
	}

	sub := claims.Sub
	if sub == "" {
		sub = claims.Username
	}
	if sub == "" {
		sub = claims.NIP
	}

	if sub == "" {
		return nil, fmt.Errorf("JWT payload missing 'sub' identifier")
	}

	nama := claims.NamaLengkap
	if nama == "" {
		nama = claims.Nama
	}
	if nama == "" {
		nama = claims.Name
	}
	if nama == "" {
		nama = claims.Username
	}

	userType := claims.UserType
	if userType == "" {
		userType = "internal"
	}

	nipLama := claims.NIPLama
	if nipLama == "" {
		nipLama = claims.NipLama
	}

	return &UserInfo{
		Sub:           sub,
		UserType:      userType,
		NamaLengkap:   nama,
		Email:         claims.Email,
		NIP:           claims.NIP,
		NIPLama:       nipLama,
		KodeSatker:    claims.KodeSatker,
		NIK:           claims.NIK,
		KodeKabupaten: claims.KodeKabupaten,
		KodeKecamatan: claims.KodeKecamatan,
		KodeDesa:      claims.KodeDesa,
		KodeSLS:       claims.KodeSLS,
		Jabatan:       claims.Jabatan,
	}, nil
}

// GetLogoutURL generates Single Log Out (SLO) redirect URL
func (c *Client) GetLogoutURL(postLogoutRedirectURI string) string {
	if postLogoutRedirectURI == "" {
		postLogoutRedirectURI = c.cfg.RedirectURI
	}

	params := url.Values{}
	params.Add("client_id", c.cfg.ClientID)
	params.Add("post_logout_redirect_uri", postLogoutRedirectURI)

	baseURL := strings.TrimRight(c.cfg.IssuerURL, "/")
	return fmt.Sprintf("%s/logout?%s", baseURL, params.Encode())
}
