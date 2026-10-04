package sso

import (
	"context"
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

// GetUserInfo retrieves user profile claims via GET /userinfo
func (c *Client) GetUserInfo(ctx context.Context, accessToken string) (*UserInfo, error) {
	userInfoURL := fmt.Sprintf("%s/userinfo", strings.TrimRight(c.cfg.IssuerURL, "/"))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, userInfoURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create userinfo request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("userinfo request failed: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read userinfo response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("userinfo endpoint returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var user UserInfo
	if err := json.Unmarshal(bodyBytes, &user); err != nil {
		return nil, fmt.Errorf("failed to decode userinfo response: %w", err)
	}

	if user.Sub == "" {
		return nil, fmt.Errorf("userinfo response missing 'sub' claim")
	}

	return &user, nil
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
