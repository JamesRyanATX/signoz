package ldappasswordauthn

import (
	"context"
	"os"
	"strings"

	"github.com/SigNoz/signoz/pkg/authn"
	"github.com/SigNoz/signoz/pkg/errors"
	"github.com/SigNoz/signoz/pkg/ldap"
	"github.com/SigNoz/signoz/pkg/modules/authdomain"
	"github.com/SigNoz/signoz/pkg/modules/user"
	"github.com/SigNoz/signoz/pkg/types"
	"github.com/SigNoz/signoz/pkg/types/authtypes"
	"github.com/SigNoz/signoz/pkg/valuer"
)

var _ authn.PasswordAuthN = (*AuthN)(nil)

type AuthN struct {
	authDomain authdomain.Module
	userSetter user.Setter
}

func New(authDomain authdomain.Module, userSetter user.Setter) *AuthN {
	return &AuthN{authDomain: authDomain, userSetter: userSetter}
}

func (a *AuthN) Authenticate(ctx context.Context, identifier string, password string, orgID valuer.UUID) (*authtypes.Identity, error) {
	ldapConfig, err := a.getLDAPConfig(ctx, identifier, orgID)
	if err != nil {
		return nil, err
	}

	ldapClient, err := ldap.NewClient(ldapConfig)
	if err != nil {
		return nil, errors.Wrapf(err, errors.TypeInternal, errors.CodeInternal, "failed to create LDAP client")
	}

	ldapAttrs, err := ldapClient.Authenticate(identifier, password)
	if err != nil {
		return nil, err
	}

	// Accept login with either email or username
	if !strings.EqualFold(ldapAttrs.Email, identifier) && !strings.EqualFold(ldapAttrs.Username, identifier) {
		return nil, errors.Newf(errors.TypeInvalidInput, errors.CodeInvalidInput, "LDAP identity mismatch: expected %s or %s", ldapAttrs.Email, ldapAttrs.Username)
	}

	ldapEmail, err := valuer.NewEmail(ldapAttrs.Email)
	if err != nil {
		return nil, errors.Wrapf(err, errors.TypeInvalidInput, errors.CodeInvalidInput, "invalid email from LDAP: %s", ldapAttrs.Email)
	}

	displayName := ldapAttrs.DisplayName
	if displayName == "" {
		displayName = ldapAttrs.Username
	}

	newUser, err := types.NewUser(displayName, ldapEmail, orgID, types.UserStatusActive)
	if err != nil {
		return nil, err
	}

	createdUser, err := a.userSetter.GetOrCreateUser(ctx, newUser)
	if err != nil {
		return nil, err
	}

	return authtypes.NewPrincipalUserIdentity(createdUser.ID, orgID, ldapEmail, authtypes.IdentNProviderTokenizer), nil
}

// getLDAPConfig finds the LDAP configuration for the given orgID.
// For username-based login (no "@"), falls back to the SIGNOZ_SSO_DEFAULT_DOMAIN env var.
func (a *AuthN) getLDAPConfig(ctx context.Context, identifier string, orgID valuer.UUID) (*authtypes.LdapConfig, error) {
	domains, err := a.authDomain.ListByOrgID(ctx, orgID)
	if err != nil {
		return nil, err
	}

	for _, domain := range domains {
		cfg := domain.AuthDomainConfig()
		if cfg.SSOEnabled && cfg.AuthNProvider == authtypes.AuthNProviderLDAP && cfg.LDAP != nil {
			return cfg.LDAP, nil
		}
	}

	// For username-based login, try the default SSO domain env var
	if !strings.Contains(identifier, "@") {
		defaultDomainName := os.Getenv("SIGNOZ_SSO_DEFAULT_DOMAIN")
		if defaultDomainName != "" {
			domain, err := a.authDomain.GetByNameAndOrgID(ctx, defaultDomainName, orgID)
			if err == nil {
				cfg := domain.AuthDomainConfig()
				if cfg.SSOEnabled && cfg.AuthNProvider == authtypes.AuthNProviderLDAP && cfg.LDAP != nil {
					return cfg.LDAP, nil
				}
			}
		}
	}

	return nil, errors.New(errors.TypeNotFound, errors.CodeNotFound, "no LDAP configuration found for this organization")
}
