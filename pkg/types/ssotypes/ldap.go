package ssotypes

import (
	"fmt"
	"strings"
)

// LdapConfig contains configuration for LDAP authentication
type LdapConfig struct {
	// Server configuration
	ServerURL  string `json:"serverUrl"`  // LDAP server URL (e.g., ldap://ldap.example.com or ldaps://ldap.example.com)
	ServerPort int    `json:"serverPort"` // LDAP server port (default: 389 for ldap, 636 for ldaps)

	// Bind configuration
	BindDN       string `json:"bindDn"`       // DN to bind with for search operations
	BindPassword string `json:"bindPassword"` // Password for bind DN (should be encrypted)

	// User search configuration
	UserBaseDN        string `json:"userBaseDn"`        // Base DN for user search (e.g., ou=users,dc=example,dc=com)
	UserFilter        string `json:"userFilter"`        // LDAP filter for user search - %s is replaced with the user's email (e.g., (mail=%s) for email or (sAMAccountName=%s) for username)
	UsernameAttribute string `json:"usernameAttribute"` // Attribute to use as username (default: uid)

	// Attribute mappings
	EmailAttribute       string `json:"emailAttribute"`       // Attribute for user email (default: mail)
	DisplayNameAttribute string `json:"displayNameAttribute"` // Attribute for display name (default: cn or displayName)

	// Group search configuration (optional)
	GroupBaseDN     string `json:"groupBaseDn,omitempty"`     // Base DN for group search
	GroupFilter     string `json:"groupFilter,omitempty"`     // LDAP filter for group search
	GroupMemberAttr string `json:"groupMemberAttr,omitempty"` // Attribute in group that contains member DNs (default: member)

	// Security settings
	UseTLS      bool `json:"useTls"`      // Use TLS connection (ldaps://)
	UseStartTLS bool `json:"useStartTls"` // Use StartTLS to upgrade connection
	SkipTLSVerify bool `json:"skipTlsVerify,omitempty"` // Skip TLS certificate verification (insecure, for testing only)

	// Advanced settings
	SearchTimeout int `json:"searchTimeout,omitempty"` // Search timeout in seconds (default: 10)
	ConnTimeout   int `json:"connTimeout,omitempty"`   // Connection timeout in seconds (default: 10)
}

// Validate checks if the LDAP configuration is valid
func (l *LdapConfig) Validate() error {
	if l.ServerURL == "" {
		return fmt.Errorf("LDAP server URL is required")
	}

	if l.ServerPort == 0 {
		// Set default ports based on protocol
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

	// Set defaults for optional fields
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

	// Validate mutual exclusivity
	if l.UseTLS && l.UseStartTLS {
		return fmt.Errorf("cannot use both TLS and StartTLS")
	}

	return nil
}

// GetServerAddress returns the full server address for connection
func (l *LdapConfig) GetServerAddress() string {
	return fmt.Sprintf("%s:%d", strings.TrimPrefix(strings.TrimPrefix(l.ServerURL, "ldap://"), "ldaps://"), l.ServerPort)
}

// IsSecure returns true if the connection uses encryption
func (l *LdapConfig) IsSecure() bool {
	return l.UseTLS || l.UseStartTLS || strings.HasPrefix(l.ServerURL, "ldaps://")
}

// LdapUserAttributes contains user attributes retrieved from LDAP
type LdapUserAttributes struct {
	Username    string
	Email       string
	DisplayName string
	Groups      []string
	DN          string // Distinguished Name of the user
}
