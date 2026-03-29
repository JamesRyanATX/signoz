package ldap

import (
	"crypto/tls"
	"fmt"
	"strings"

	"github.com/SigNoz/signoz/pkg/errors"
	"github.com/SigNoz/signoz/pkg/types/ssotypes"
	"github.com/go-ldap/ldap/v3"
)

// Client is an LDAP client wrapper for authentication operations
type Client struct {
	config *ssotypes.LdapConfig
}

// NewClient creates a new LDAP client with the given configuration
func NewClient(config *ssotypes.LdapConfig) (*Client, error) {
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid LDAP configuration: %w", err)
	}

	return &Client{
		config: config,
	}, nil
}

// connect establishes a connection to the LDAP server
func (c *Client) connect() (*ldap.Conn, error) {
	serverAddr := c.config.GetServerAddress()

	var conn *ldap.Conn
	var err error

	// Determine connection type
	if c.config.UseTLS || strings.HasPrefix(c.config.ServerURL, "ldaps://") {
		// Use TLS connection
		tlsConfig := &tls.Config{
			ServerName:         strings.TrimPrefix(strings.TrimPrefix(c.config.ServerURL, "ldap://"), "ldaps://"),
			InsecureSkipVerify: c.config.SkipTLSVerify,
		}
		conn, err = ldap.DialTLS("tcp", serverAddr, tlsConfig)
	} else {
		// Use plain connection
		conn, err = ldap.Dial("tcp", serverAddr)
	}

	if err != nil {
		return nil, errors.Wrapf(err, errors.TypeInternal, errors.CodeInternal, "failed to connect to LDAP server at %s", serverAddr)
	}

	// Upgrade to TLS if StartTLS is enabled
	if c.config.UseStartTLS && !c.config.UseTLS {
		tlsConfig := &tls.Config{
			ServerName:         strings.TrimPrefix(strings.TrimPrefix(c.config.ServerURL, "ldap://"), "ldaps://"),
			InsecureSkipVerify: c.config.SkipTLSVerify,
		}
		if err := conn.StartTLS(tlsConfig); err != nil {
			conn.Close()
			return nil, errors.Wrapf(err, errors.TypeInternal, errors.CodeInternal, "failed to start TLS")
		}
	}

	return conn, nil
}

// bind performs a bind operation with the given DN and password
func (c *Client) bind(conn *ldap.Conn, dn, password string) error {
	if err := conn.Bind(dn, password); err != nil {
		if ldap.IsErrorWithCode(err, ldap.LDAPResultInvalidCredentials) {
			return errors.New(errors.TypeInvalidInput, errors.CodeInvalidInput, "invalid LDAP credentials")
		}
		return errors.Wrapf(err, errors.TypeInternal, errors.CodeInternal, "LDAP bind failed for %s", dn)
	}
	return nil
}

// searchUser searches for a user by username and returns the user DN and attributes
func (c *Client) searchUser(conn *ldap.Conn, username string) (*ldap.Entry, error) {
	// Normalize to lowercase for case-insensitive matching
	// Most LDAP servers handle case-insensitivity, but this ensures consistent behavior
	normalizedUsername := strings.ToLower(username)

	// Build search filter - replace %s with username
	searchFilter := strings.ReplaceAll(c.config.UserFilter, "%s", ldap.EscapeFilter(normalizedUsername))

	// Build attribute list to retrieve
	attributes := []string{
		c.config.EmailAttribute,
		c.config.DisplayNameAttribute,
		c.config.UsernameAttribute,
	}

	// Create search request
	searchRequest := ldap.NewSearchRequest(
		c.config.UserBaseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		1, // Size limit - we only need one user
		c.config.SearchTimeout,
		false,
		searchFilter,
		attributes,
		nil,
	)

	// Perform search
	sr, err := conn.Search(searchRequest)
	if err != nil {
		return nil, errors.Wrapf(err, errors.TypeInternal, errors.CodeInternal, "LDAP search failed for user %s", username)
	}

	if len(sr.Entries) == 0 {
		return nil, errors.Newf(errors.TypeNotFound, errors.CodeNotFound, "user %s not found in LDAP", username)
	}

	if len(sr.Entries) > 1 {
		return nil, errors.Newf(errors.TypeInternal, errors.CodeInternal, "multiple users found for %s in LDAP", username)
	}

	return sr.Entries[0], nil
}

// searchUserGroups searches for groups that a user belongs to
func (c *Client) searchUserGroups(conn *ldap.Conn, userDN string) ([]string, error) {
	// If group search is not configured, return empty list
	if c.config.GroupBaseDN == "" || c.config.GroupFilter == "" {
		return []string{}, nil
	}

	// Build group search filter - replace %s with user DN
	searchFilter := strings.ReplaceAll(c.config.GroupFilter, "%s", ldap.EscapeFilter(userDN))

	// Create search request for groups
	searchRequest := ldap.NewSearchRequest(
		c.config.GroupBaseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0, // No size limit
		c.config.SearchTimeout,
		false,
		searchFilter,
		[]string{"cn"}, // Get group common name
		nil,
	)

	// Perform search
	sr, err := conn.Search(searchRequest)
	if err != nil {
		// Don't fail authentication if group search fails, just log and return empty
		return []string{}, nil
	}

	groups := make([]string, 0, len(sr.Entries))
	for _, entry := range sr.Entries {
		if cn := entry.GetAttributeValue("cn"); cn != "" {
			groups = append(groups, cn)
		}
	}

	return groups, nil
}

// extractUserAttributes extracts user attributes from LDAP entry
func (c *Client) extractUserAttributes(entry *ldap.Entry, groups []string) *ssotypes.LdapUserAttributes {
	return &ssotypes.LdapUserAttributes{
		Username:    entry.GetAttributeValue(c.config.UsernameAttribute),
		Email:       entry.GetAttributeValue(c.config.EmailAttribute),
		DisplayName: entry.GetAttributeValue(c.config.DisplayNameAttribute),
		DN:          entry.DN,
		Groups:      groups,
	}
}

// Authenticate authenticates a user against LDAP using their email address and returns user attributes
// SigNoz uses email-based login, so the identifier parameter should be the user's email address
func (c *Client) Authenticate(emailOrIdentifier, password string) (*ssotypes.LdapUserAttributes, error) {
	if emailOrIdentifier == "" || password == "" {
		return nil, errors.New(errors.TypeInvalidInput, errors.CodeInvalidInput, "email and password are required")
	}

	// Connect to LDAP server
	conn, err := c.connect()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// Bind with service account to search for user
	if err := c.bind(conn, c.config.BindDN, c.config.BindPassword); err != nil {
		return nil, fmt.Errorf("failed to bind with service account: %w", err)
	}

	// Search for user using the configured UserFilter
	// The admin configures this to search by whatever field makes sense (email, username, etc.)
	userEntry, err := c.searchUser(conn, emailOrIdentifier)
	if err != nil {
		return nil, err
	}

	userDN := userEntry.DN

	// Bind as the user to verify credentials
	if err := c.bind(conn, userDN, password); err != nil {
		return nil, err
	}

	// Re-bind with service account to search for groups
	if err := c.bind(conn, c.config.BindDN, c.config.BindPassword); err != nil {
		// If group search fails, continue without groups
		groups, _ := c.searchUserGroups(conn, userDN)
		return c.extractUserAttributes(userEntry, groups), nil
	}

	// Search for user groups
	groups, err := c.searchUserGroups(conn, userDN)
	if err != nil {
		// If group search fails, continue without groups
		return c.extractUserAttributes(userEntry, []string{}), nil
	}

	return c.extractUserAttributes(userEntry, groups), nil
}

// TestConnection tests the connection to the LDAP server and validates configuration
func (c *Client) TestConnection() error {
	// Connect to LDAP server
	conn, err := c.connect()
	if err != nil {
		return err
	}
	defer conn.Close()

	// Test bind with service account
	if err := c.bind(conn, c.config.BindDN, c.config.BindPassword); err != nil {
		return err
	}

	// Test if user base DN is accessible
	searchRequest := ldap.NewSearchRequest(
		c.config.UserBaseDN,
		ldap.ScopeBaseObject, // Just check if base exists
		ldap.NeverDerefAliases,
		1,
		c.config.SearchTimeout,
		false,
		"(objectClass=*)",
		[]string{"dn"},
		nil,
	)

	_, err = conn.Search(searchRequest)
	if err != nil {
		return errors.Wrapf(err, errors.TypeInternal, errors.CodeInternal, "user base DN %s is not accessible", c.config.UserBaseDN)
	}

	return nil
}
