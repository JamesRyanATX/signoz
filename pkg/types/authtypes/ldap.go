package authtypes

import (
	"fmt"
	"strings"
)

// LdapConfig contains configuration for LDAP authentication
type LdapConfig struct {
	// Server configuration
	ServerURL  string `json:"serverUrl"`
	ServerPort int    `json:"serverPort"`

	// Bind configuration
	BindDN       string `json:"bindDn"`
	BindPassword string `json:"bindPassword"`

	// User search configuration
	UserBaseDN        string `json:"userBaseDn"`
	UserFilter        string `json:"userFilter"`
	UsernameAttribute string `json:"usernameAttribute"`

	// Attribute mappings
	EmailAttribute       string `json:"emailAttribute"`
	DisplayNameAttribute string `json:"displayNameAttribute"`

	// Group search configuration (optional)
	GroupBaseDN     string `json:"groupBaseDn,omitempty"`
	GroupFilter     string `json:"groupFilter,omitempty"`
	GroupMemberAttr string `json:"groupMemberAttr,omitempty"`

	// Security settings
	UseTLS        bool `json:"useTls"`
	UseStartTLS   bool `json:"useStartTls"`
	SkipTLSVerify bool `json:"skipTlsVerify,omitempty"`

	// Advanced settings
	SearchTimeout int `json:"searchTimeout,omitempty"`
	ConnTimeout   int `json:"connTimeout,omitempty"`
}

// Validate checks if the LDAP configuration is valid and fills in defaults.
func (l *LdapConfig) Validate() error {
	if l.ServerURL == "" {
		return fmt.Errorf("LDAP server URL is required")
	}

	if l.ServerPort == 0 {
		if strings.HasPrefix(l.ServerURL, "ldaps://") || l.UseTLS {
			l.ServerPort = 636
		} else {
			l.ServerPort = 389
		}
	}

	if l.BindDN == "" {
		return fmt.Errorf("LDAP bind DN is required")
	}

	if l.BindPassword == "" {
		return fmt.Errorf("LDAP bind password is required")
	}

	if l.UserBaseDN == "" {
		return fmt.Errorf("LDAP user base DN is required")
	}

	if l.UserFilter == "" {
		return fmt.Errorf("LDAP user filter is required")
	}

	if l.UsernameAttribute == "" {
		l.UsernameAttribute = "uid"
	}

	if l.EmailAttribute == "" {
		l.EmailAttribute = "mail"
	}

	if l.DisplayNameAttribute == "" {
		l.DisplayNameAttribute = "cn"
	}

	if l.GroupMemberAttr == "" {
		l.GroupMemberAttr = "member"
	}

	if l.SearchTimeout == 0 {
		l.SearchTimeout = 10
	}

	if l.ConnTimeout == 0 {
		l.ConnTimeout = 10
	}

	if l.UseTLS && l.UseStartTLS {
		return fmt.Errorf("cannot use both TLS and StartTLS")
	}

	return nil
}

// GetServerAddress returns the host:port for dialing.
func (l *LdapConfig) GetServerAddress() string {
	return fmt.Sprintf("%s:%d", strings.TrimPrefix(strings.TrimPrefix(l.ServerURL, "ldap://"), "ldaps://"), l.ServerPort)
}

// IsSecure returns true if the connection uses TLS.
func (l *LdapConfig) IsSecure() bool {
	return l.UseTLS || l.UseStartTLS || strings.HasPrefix(l.ServerURL, "ldaps://")
}

// LdapUserAttributes contains user attributes retrieved from LDAP.
type LdapUserAttributes struct {
	Username    string
	Email       string
	DisplayName string
	Groups      []string
	DN          string
}
