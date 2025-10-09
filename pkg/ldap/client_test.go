package ldap

import (
	"testing"

	"github.com/SigNoz/signoz/pkg/types/ssotypes"
)

func TestLdapConfigValidation(t *testing.T) {
	tests := []struct {
		name      string
		config    *ssotypes.LdapConfig
		wantError bool
	}{
		{
			name: "valid LDAP configuration",
			config: &ssotypes.LdapConfig{
				ServerURL:            "ldap://ldap.example.com",
				ServerPort:           389,
				BindDN:               "cn=admin,dc=example,dc=com",
				BindPassword:         "password",
				UserBaseDN:           "ou=users,dc=example,dc=com",
				UserFilter:           "(uid=%s)",
				UsernameAttribute:    "uid",
				EmailAttribute:       "mail",
				DisplayNameAttribute: "cn",
			},
			wantError: false,
		},
		{
			name: "missing server URL",
			config: &ssotypes.LdapConfig{
				ServerPort:   389,
				BindDN:       "cn=admin,dc=example,dc=com",
				BindPassword: "password",
				UserBaseDN:   "ou=users,dc=example,dc=com",
				UserFilter:   "(uid=%s)",
			},
			wantError: true,
		},
		{
			name: "missing bind DN",
			config: &ssotypes.LdapConfig{
				ServerURL:    "ldap://ldap.example.com",
				ServerPort:   389,
				BindPassword: "password",
				UserBaseDN:   "ou=users,dc=example,dc=com",
				UserFilter:   "(uid=%s)",
			},
			wantError: true,
		},
		{
			name: "missing user base DN",
			config: &ssotypes.LdapConfig{
				ServerURL:    "ldap://ldap.example.com",
				ServerPort:   389,
				BindDN:       "cn=admin,dc=example,dc=com",
				BindPassword: "password",
				UserFilter:   "(uid=%s)",
			},
			wantError: true,
		},
		{
			name: "both TLS and StartTLS enabled",
			config: &ssotypes.LdapConfig{
				ServerURL:    "ldap://ldap.example.com",
				ServerPort:   389,
				BindDN:       "cn=admin,dc=example,dc=com",
				BindPassword: "password",
				UserBaseDN:   "ou=users,dc=example,dc=com",
				UserFilter:   "(uid=%s)",
				UseTLS:       true,
				UseStartTLS:  true,
			},
			wantError: true,
		},
		{
			name: "auto-set port for ldaps",
			config: &ssotypes.LdapConfig{
				ServerURL:    "ldaps://ldap.example.com",
				BindDN:       "cn=admin,dc=example,dc=com",
				BindPassword: "password",
				UserBaseDN:   "ou=users,dc=example,dc=com",
				UserFilter:   "(uid=%s)",
			},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantError {
				t.Errorf("Validate() error = %v, wantError %v", err, tt.wantError)
			}

			// Check that defaults are set when validation succeeds
			if err == nil {
				if tt.config.UsernameAttribute == "" {
					t.Error("UsernameAttribute should have default value")
				}
				if tt.config.EmailAttribute == "" {
					t.Error("EmailAttribute should have default value")
				}
				if tt.config.DisplayNameAttribute == "" {
					t.Error("DisplayNameAttribute should have default value")
				}
				if tt.config.SearchTimeout == 0 {
					t.Error("SearchTimeout should have default value")
				}
				if tt.config.ConnTimeout == 0 {
					t.Error("ConnTimeout should have default value")
				}
			}
		})
	}
}

func TestGetServerAddress(t *testing.T) {
	tests := []struct {
		name     string
		config   *ssotypes.LdapConfig
		expected string
	}{
		{
			name: "ldap with standard port",
			config: &ssotypes.LdapConfig{
				ServerURL:  "ldap://ldap.example.com",
				ServerPort: 389,
			},
			expected: "ldap.example.com:389",
		},
		{
			name: "ldaps with standard port",
			config: &ssotypes.LdapConfig{
				ServerURL:  "ldaps://ldap.example.com",
				ServerPort: 636,
			},
			expected: "ldap.example.com:636",
		},
		{
			name: "custom port",
			config: &ssotypes.LdapConfig{
				ServerURL:  "ldap://ldap.example.com",
				ServerPort: 1389,
			},
			expected: "ldap.example.com:1389",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.config.GetServerAddress()
			if got != tt.expected {
				t.Errorf("GetServerAddress() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestIsSecure(t *testing.T) {
	tests := []struct {
		name     string
		config   *ssotypes.LdapConfig
		expected bool
	}{
		{
			name: "ldaps protocol",
			config: &ssotypes.LdapConfig{
				ServerURL: "ldaps://ldap.example.com",
			},
			expected: true,
		},
		{
			name: "UseTLS enabled",
			config: &ssotypes.LdapConfig{
				ServerURL: "ldap://ldap.example.com",
				UseTLS:    true,
			},
			expected: true,
		},
		{
			name: "UseStartTLS enabled",
			config: &ssotypes.LdapConfig{
				ServerURL:   "ldap://ldap.example.com",
				UseStartTLS: true,
			},
			expected: true,
		},
		{
			name: "plain ldap",
			config: &ssotypes.LdapConfig{
				ServerURL: "ldap://ldap.example.com",
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.config.IsSecure()
			if got != tt.expected {
				t.Errorf("IsSecure() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestNewClient(t *testing.T) {
	tests := []struct {
		name      string
		config    *ssotypes.LdapConfig
		wantError bool
	}{
		{
			name: "valid config",
			config: &ssotypes.LdapConfig{
				ServerURL:    "ldap://ldap.example.com",
				ServerPort:   389,
				BindDN:       "cn=admin,dc=example,dc=com",
				BindPassword: "password",
				UserBaseDN:   "ou=users,dc=example,dc=com",
				UserFilter:   "(uid=%s)",
			},
			wantError: false,
		},
		{
			name: "invalid config - missing fields",
			config: &ssotypes.LdapConfig{
				ServerURL: "ldap://ldap.example.com",
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.config)
			if (err != nil) != tt.wantError {
				t.Errorf("NewClient() error = %v, wantError %v", err, tt.wantError)
			}
			if !tt.wantError && client == nil {
				t.Error("NewClient() returned nil client for valid config")
			}
		})
	}
}
