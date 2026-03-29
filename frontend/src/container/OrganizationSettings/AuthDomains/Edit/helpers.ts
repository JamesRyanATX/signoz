import {
	AuthDomain,
	GOOGLE_AUTH,
	GoogleAuthConfig,
	isGoogleAuthConfig,
	isLDAPConfig,
	isSAMLConfig,
	LDAP,
	LDAPConfig,
	SAML,
	SAMLConfig,
} from 'types/api/SAML/listDomain';

export function parseSamlForm(
	current: AuthDomain,
	formValues: AuthDomain,
): SAMLConfig | undefined {
	if (current?.ssoType === SAML && isSAMLConfig(formValues?.samlConfig)) {
		return {
			...current.samlConfig,
			...formValues?.samlConfig,
		};
	}

	return current.samlConfig;
}

export function parseGoogleAuthForm(
	current: AuthDomain,
	formValues: AuthDomain,
): GoogleAuthConfig | undefined {
	if (
		current?.ssoType === GOOGLE_AUTH &&
		isGoogleAuthConfig(formValues?.googleAuthConfig)
	) {
		return {
			...current.googleAuthConfig,
			...formValues?.googleAuthConfig,
		};
	}

	return current.googleAuthConfig;
}

export function parseLdapForm(
	current: AuthDomain,
	formValues: AuthDomain,
): LDAPConfig | undefined {
	if (current?.ssoType === LDAP && isLDAPConfig(formValues?.ldapConfig)) {
		return {
			...current.ldapConfig,
			...formValues?.ldapConfig,
		};
	}

	return current.ldapConfig;
}
