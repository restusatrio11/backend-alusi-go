package sso_test

import (
	"strings"
	"testing"

	"backend-alusi-go/config"
	"backend-alusi-go/pkg/sso"
)

func TestSSOClient_URLs(t *testing.T) {
	cfg := &config.SSOConfig{
		ClientID:     "app_1791057531231_943f9ac1",
		ClientSecret: "test_secret",
		RedirectURI:  "https://aron.bps.web.id/callback",
		IssuerURL:    "https://otp-dev.bps.web.id",
	}

	client := sso.NewClient(cfg)

	// 1. Test Auth URL
	authURL := client.GetAuthURL("state_abc_123")
	if !strings.HasPrefix(authURL, "https://otp-dev.bps.web.id/authorize?") {
		t.Errorf("Unexpected auth URL prefix: %s", authURL)
	}
	if !strings.Contains(authURL, "client_id=app_1791057531231_943f9ac1") {
		t.Errorf("Auth URL missing client_id: %s", authURL)
	}
	if !strings.Contains(authURL, "state=state_abc_123") {
		t.Errorf("Auth URL missing state: %s", authURL)
	}
	if !strings.Contains(authURL, "response_type=code") {
		t.Errorf("Auth URL missing response_type: %s", authURL)
	}

	// 2. Test Logout URL
	logoutURL := client.GetLogoutURL("https://aron.bps.web.id/callback")
	if !strings.HasPrefix(logoutURL, "https://otp-dev.bps.web.id/logout?") {
		t.Errorf("Unexpected logout URL prefix: %s", logoutURL)
	}
	if !strings.Contains(logoutURL, "post_logout_redirect_uri=https%3A%2F%2Faron.bps.web.id%2Fcallback") {
		t.Errorf("Logout URL missing post_logout_redirect_uri: %s", logoutURL)
	}
}
