package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	"backend-alusi-go/internal/domain"
	"backend-alusi-go/internal/repository/postgres"
	"backend-alusi-go/pkg/crypto"
	"backend-alusi-go/pkg/jwt"
	"backend-alusi-go/pkg/sso"

	"github.com/rs/zerolog/log"
)

type ManualLoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AuthUsecase struct {
	userRepo   *postgres.UserRepo
	ssoClient  *sso.Client
	jwtService *jwt.JWTService
}

func NewAuthUsecase(userRepo *postgres.UserRepo, ssoClient *sso.Client, jwtService *jwt.JWTService) *AuthUsecase {
	return &AuthUsecase{
		userRepo:   userRepo,
		ssoClient:  ssoClient,
		jwtService: jwtService,
	}
}

// GetLoginURL generates the SSO authorization redirect URL
func (u *AuthUsecase) GetLoginURL(state string) string {
	return u.ssoClient.GetAuthURL(state)
}

// HandleSSOCallback processes the SSO callback, performs JIT provisioning, and returns user and session token
func (u *AuthUsecase) HandleSSOCallback(ctx context.Context, code string) (*domain.User, string, error) {
	// 1. Exchange authorization code for SSO tokens
	tok, err := u.ssoClient.ExchangeCode(ctx, code)
	if err != nil {
		return nil, "", fmt.Errorf("failed to exchange auth code: %w", err)
	}

	// 2. Fetch User Claims from /userinfo
	userInfo, err := u.ssoClient.GetUserInfo(ctx, tok.AccessToken)
	if err != nil {
		return nil, "", fmt.Errorf("failed to fetch userinfo from SSO: %w", err)
	}

	// 3. Resolve Satker
	satkerKode := userInfo.KodeSatker
	if satkerKode == "" && userInfo.KodeKabupaten != "" {
		satkerKode = userInfo.KodeKabupaten
	}

	var satkerID *int
	if satkerKode != "" {
		satker, err := u.userRepo.GetSatkerByKode(ctx, satkerKode)
		if err == nil && satker != nil {
			satkerID = &satker.ID
		}
	}

	// 4. Just-In-Time (JIT) Provisioning
	user, err := u.userRepo.GetBySSOSub(ctx, userInfo.Sub)
	if err != nil {
		return nil, "", fmt.Errorf("database query error: %w", err)
	}

	now := time.Now()
	metadata := map[string]interface{}{}
	if userInfo.UserType == "external" {
		metadata["kode_kabupaten"] = userInfo.KodeKabupaten
		metadata["kode_kecamatan"] = userInfo.KodeKecamatan
		metadata["kode_desa"] = userInfo.KodeDesa
		metadata["kode_sls"] = userInfo.KodeSLS
		metadata["jabatan"] = userInfo.Jabatan
	}

	var nip *string
	if userInfo.NIP != "" {
		nip = &userInfo.NIP
	}

	var nik *string
	if userInfo.NIK != "" {
		nik = &userInfo.NIK
	}

	if user == nil {
		// New User -> Create
		newUser := &domain.User{
			SSOSub:      userInfo.Sub,
			UserType:    userInfo.UserType,
			NIP:         nip,
			NIK:         nik,
			Nama:        userInfo.NamaLengkap,
			Email:       userInfo.Email,
			SatkerID:    satkerID,
			Status:      "active",
			Metadata:    metadata,
			LastLoginAt: &now,
		}

		if err := u.userRepo.Create(ctx, newUser); err != nil {
			return nil, "", fmt.Errorf("failed to create user in database: %w", err)
		}
		user = newUser

		// Assign default role based on user_type
		var defaultRoleName string
		if userInfo.UserType == "internal" {
			defaultRoleName = "pegawai"
		} else {
			if strings.Contains(strings.ToLower(userInfo.Jabatan), "lapangan") || strings.Contains(strings.ToLower(userInfo.Jabatan), "ppl") {
				defaultRoleName = "petugas_lapangan"
			} else {
				defaultRoleName = "mitra"
			}
		}

		role, err := u.userRepo.GetRoleByNama(ctx, defaultRoleName)
		if err == nil && role != nil {
			_ = u.userRepo.AssignRoles(ctx, user.ID, []int{role.ID})
			user.Roles = []domain.Role{*role}
		}

		log.Info().
			Str("sso_sub", user.SSOSub).
			Str("user_type", user.UserType).
			Str("nama", user.Nama).
			Msg("New user provisioned via JIT SSO")
	} else {
		// Existing User -> Synchronize profile & update login time
		user.UserType = userInfo.UserType
		user.Nama = userInfo.NamaLengkap
		user.Email = userInfo.Email
		if nip != nil {
			user.NIP = nip
		}
		if nik != nil {
			user.NIK = nik
		}
		if satkerID != nil {
			user.SatkerID = satkerID
		}
		user.Metadata = metadata

		if err := u.userRepo.Update(ctx, user); err != nil {
			log.Warn().Err(err).Msg("Failed to sync updated user profile")
		}
		_ = u.userRepo.UpdateLastLogin(ctx, user.ID, now)
	}

	// Reload full user with satker & roles
	fullUser, err := u.userRepo.GetByID(ctx, user.ID)
	if err == nil && fullUser != nil {
		user = fullUser
	}

	// 5. Generate Portal Session JWT
	resolvedSatkerKode := ""
	if user.Satker != nil {
		resolvedSatkerKode = user.Satker.Kode
	} else if satkerKode != "" {
		resolvedSatkerKode = satkerKode
	}

	sessionToken, _, err := u.jwtService.GenerateSessionToken(user, resolvedSatkerKode)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate session token: %w", err)
	}

	return user, sessionToken, nil
}

// GetLogoutURL returns Single Log Out URL
func (u *AuthUsecase) GetLogoutURL(postLogoutRedirectURI string) string {
	return u.ssoClient.GetLogoutURL(postLogoutRedirectURI)
}

// GetProfile returns authenticated user details
func (u *AuthUsecase) GetProfile(ctx context.Context, userID int) (*domain.User, error) {
	if u.userRepo == nil {
		return nil, fmt.Errorf("user repository is not initialized")
	}
	return u.userRepo.GetByID(ctx, userID)
}

// ManualLogin verifies username/email/NIP and password, updates last login, and returns user and session token
func (u *AuthUsecase) ManualLogin(ctx context.Context, input ManualLoginInput) (*domain.User, string, error) {
	if u.userRepo == nil {
		return nil, "", fmt.Errorf("user repository is not initialized")
	}

	trimmedUsername := strings.TrimSpace(input.Username)
	if trimmedUsername == "" || input.Password == "" {
		return nil, "", fmt.Errorf("username dan password wajib diisi")
	}

	user, err := u.userRepo.GetByUsernameOrEmailOrNIP(ctx, trimmedUsername)
	if err != nil {
		return nil, "", fmt.Errorf("gagal memeriksa kredensial pengguna: %w", err)
	}

	if user == nil {
		return nil, "", fmt.Errorf("username atau password salah")
	}

	if user.Status != "active" {
		return nil, "", fmt.Errorf("akun Anda berstatus non-aktif, silakan hubungi administrator")
	}

	if user.PasswordHash == nil || *user.PasswordHash == "" {
		return nil, "", fmt.Errorf("akun ini dikonfigurasi menggunakan SSO BPS, silakan login dengan SSO")
	}

	if !crypto.CheckPasswordHash(input.Password, *user.PasswordHash) {
		return nil, "", fmt.Errorf("username atau password salah")
	}

	// Update last login timestamp
	now := time.Now()
	_ = u.userRepo.UpdateLastLogin(ctx, user.ID, now)
	user.LastLoginAt = &now

	// Reload full user with satker & roles
	fullUser, err := u.userRepo.GetByID(ctx, user.ID)
	if err == nil && fullUser != nil {
		user = fullUser
	}

	// Generate Portal Session JWT
	resolvedSatkerKode := ""
	if user.Satker != nil {
		resolvedSatkerKode = user.Satker.Kode
	}

	sessionToken, _, err := u.jwtService.GenerateSessionToken(user, resolvedSatkerKode)
	if err != nil {
		return nil, "", fmt.Errorf("gagal membuat token sesi: %w", err)
	}

	return user, sessionToken, nil
}

