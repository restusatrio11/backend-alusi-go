package config

import (
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

	idleTime, err := time.ParseDuration(v.GetString("DB_MAX_IDLE_TIME"))
	if err != nil {
		idleTime = 15 * time.Minute
	}

	originsRaw := v.GetString("CORS_ALLOWED_ORIGINS")
	origins := strings.Split(originsRaw, ",")
	for i := range origins {
		origins[i] = strings.TrimSpace(origins[i])
	}

	cfg := &Config{
		App: AppConfig{
			Name:  v.GetString("APP_NAME"),
			Env:   v.GetString("APP_ENV"),
			Port:  v.GetString("APP_PORT"),
			Debug: v.GetBool("APP_DEBUG"),
		},
		Database: DatabaseConfig{
			Host:         v.GetString("DB_HOST"),
			Port:         v.GetString("DB_PORT"),
			User:         v.GetString("DB_USER"),
			Password:     v.GetString("DB_PASSWORD"),
			Name:         v.GetString("DB_NAME"),
			SSLMode:      v.GetString("DB_SSLMODE"),
			MaxOpenConns: v.GetInt("DB_MAX_OPEN_CONNS"),
			MaxIdleConns: v.GetInt("DB_MAX_IDLE_CONNS"),
			MaxIdleTime:  idleTime,
		},
		SSO: SSOConfig{
			ClientID:     v.GetString("SSO_CLIENT_ID"),
			ClientSecret: v.GetString("SSO_CLIENT_SECRET"),
			RedirectURI:  v.GetString("SSO_REDIRECT_URI"),
			IssuerURL:    v.GetString("SSO_ISSUER_URL"),
			JWKSURI:      v.GetString("SSO_JWKS_URI"),
		},
		JWT: JWTConfig{
			Secret:          v.GetString("JWT_SECRET"),
			ExpirationHours: v.GetInt("JWT_EXPIRATION_HOURS"),
			CookieDomain:    v.GetString("COOKIE_DOMAIN"),
			CookieSecure:    v.GetBool("COOKIE_SECURE"),
		},
		CORS: CORSConfig{
			AllowedOrigins: origins,
		},
	}

	return cfg, nil
}
