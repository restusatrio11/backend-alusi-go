package config

import (
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig
	Database DatabaseConfig
	SSO      SSOConfig
	JWT      JWTConfig
	CORS     CORSConfig
}

type AppConfig struct {
	Name  string `mapstructure:"name"`
	Env   string `mapstructure:"env"`
	Port  string `mapstructure:"port"`
	Debug bool   `mapstructure:"debug"`
}

type DatabaseConfig struct {
	DSN          string        `mapstructure:"dsn"`
	Host         string        `mapstructure:"host"`
	Port         string        `mapstructure:"port"`
	User         string        `mapstructure:"user"`
	Password     string        `mapstructure:"password"`
	Name         string        `mapstructure:"name"`
	SSLMode      string        `mapstructure:"sslmode"`
	MaxOpenConns int           `mapstructure:"max_open_conns"`
	MaxIdleConns int           `mapstructure:"max_idle_conns"`
	MaxIdleTime  time.Duration `mapstructure:"max_idle_time"`
}

type SSOConfig struct {
	ClientID     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
	RedirectURI  string `mapstructure:"redirect_uri"`
	IssuerURL    string `mapstructure:"issuer_url"`
	JWKSURI      string `mapstructure:"jwks_uri"`
}

type JWTConfig struct {
	Secret          string `mapstructure:"secret"`
	ExpirationHours int    `mapstructure:"expiration_hours"`
	CookieDomain    string `mapstructure:"cookie_domain"`
	CookieSecure    bool   `mapstructure:"cookie_secure"`
}

type CORSConfig struct {
	AllowedOrigins []string
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

// LoadConfig loads application configuration from .env or environment variables
func LoadConfig() (*Config, error) {
	v := viper.New()
	v.SetConfigFile(".env")
	v.SetConfigType("env")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Default values
	v.SetDefault("APP_NAME", "portal-bps-sumut-backend")
	v.SetDefault("APP_ENV", "development")
	v.SetDefault("APP_PORT", "8080")
	v.SetDefault("APP_DEBUG", true)

	v.SetDefault("DATABASE_URL", "")
	v.SetDefault("DB_DSN", "")
	v.SetDefault("DB_HOST", "localhost")
	v.SetDefault("DB_PORT", "5432")
	v.SetDefault("DB_USER", "postgres")
	v.SetDefault("DB_PASSWORD", "postgres")
	v.SetDefault("DB_NAME", "alusi_portal_db")
	v.SetDefault("DB_SSLMODE", "disable")
	v.SetDefault("DB_MAX_OPEN_CONNS", 25)
	v.SetDefault("DB_MAX_IDLE_CONNS", 10)
	v.SetDefault("DB_MAX_IDLE_TIME", "15m")

	v.SetDefault("JWT_EXPIRATION_HOURS", 24)
	v.SetDefault("COOKIE_DOMAIN", "localhost")
	v.SetDefault("COOKIE_SECURE", false)
	v.SetDefault("CORS_ALLOWED_ORIGINS", "http://localhost:3000,http://localhost:5173")

	_ = v.ReadInConfig()

	idleTimeStr := getEnv("DB_MAX_IDLE_TIME", v.GetString("DB_MAX_IDLE_TIME"))
	idleTime, err := time.ParseDuration(idleTimeStr)
	if err != nil {
		idleTime = 15 * time.Minute
	}

	originsRaw := getEnv("CORS_ALLOWED_ORIGINS", v.GetString("CORS_ALLOWED_ORIGINS"))
	origins := strings.Split(originsRaw, ",")
	for i := range origins {
		origins[i] = strings.TrimSpace(origins[i])
	}

	dbDSN := os.Getenv("DATABASE_URL")
	if dbDSN == "" {
		dbDSN = os.Getenv("DB_DSN")
	}
	if dbDSN == "" {
		dbDSN = v.GetString("DATABASE_URL")
	}
	if dbDSN == "" {
		dbDSN = v.GetString("DB_DSN")
	}

	cfg := &Config{
		App: AppConfig{
			Name:  getEnv("APP_NAME", v.GetString("APP_NAME")),
			Env:   getEnv("APP_ENV", v.GetString("APP_ENV")),
			Port:  getEnv("APP_PORT", v.GetString("APP_PORT")),
			Debug: v.GetBool("APP_DEBUG"),
		},
		Database: DatabaseConfig{
			DSN:          dbDSN,
			Host:         getEnv("DB_HOST", v.GetString("DB_HOST")),
			Port:         getEnv("DB_PORT", v.GetString("DB_PORT")),
			User:         getEnv("DB_USER", v.GetString("DB_USER")),
			Password:     getEnv("DB_PASSWORD", v.GetString("DB_PASSWORD")),
			Name:         getEnv("DB_NAME", v.GetString("DB_NAME")),
			SSLMode:      getEnv("DB_SSLMODE", v.GetString("DB_SSLMODE")),
			MaxOpenConns: v.GetInt("DB_MAX_OPEN_CONNS"),
			MaxIdleConns: v.GetInt("DB_MAX_IDLE_CONNS"),
			MaxIdleTime:  idleTime,
		},
		SSO: SSOConfig{
			ClientID:     getEnv("SSO_CLIENT_ID", v.GetString("SSO_CLIENT_ID")),
			ClientSecret: getEnv("SSO_CLIENT_SECRET", v.GetString("SSO_CLIENT_SECRET")),
			RedirectURI:  getEnv("SSO_REDIRECT_URI", v.GetString("SSO_REDIRECT_URI")),
			IssuerURL:    getEnv("SSO_ISSUER_URL", v.GetString("SSO_ISSUER_URL")),
			JWKSURI:      getEnv("SSO_JWKS_URI", v.GetString("SSO_JWKS_URI")),
		},
		JWT: JWTConfig{
			Secret:          getEnv("JWT_SECRET", v.GetString("JWT_SECRET")),
			ExpirationHours: v.GetInt("JWT_EXPIRATION_HOURS"),
			CookieDomain:    getEnv("COOKIE_DOMAIN", v.GetString("COOKIE_DOMAIN")),
			CookieSecure:    v.GetBool("COOKIE_SECURE"),
		},
		CORS: CORSConfig{
			AllowedOrigins: origins,
		},
	}

	return cfg, nil
}
