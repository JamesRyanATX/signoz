import { AuthDomain, GOOGLE_AUTH, LDAP, SAML } from 'types/api/SAML/listDomain';

export const ConfigureSsoButtonText = (
	ssoType: AuthDomain['ssoType'],
): string => {
	switch (ssoType) {
		case SAML:
			return 'Edit SAML';
		case GOOGLE_AUTH:
			return 'Edit Google Auth';
		case LDAP:
			return 'Edit LDAP';
		default:
			return 'Configure SSO';
	}
};

export const EditModalTitleText = (
	ssoType: AuthDomain['ssoType'] | undefined,
): string => {
	switch (ssoType) {
		case SAML:
			return 'Edit SAML Configuration';
		case GOOGLE_AUTH:
			return 'Edit Google Authentication';
		case LDAP:
			return 'Edit LDAP Configuration';
		default:
			return 'Configure SSO';
	}
};

export const isSSOConfigValid = (domain: AuthDomain): boolean => {
	switch (domain.ssoType) {
		case SAML:
			return (
				domain.samlConfig?.samlCert?.length !== 0 &&
				domain.samlConfig?.samlEntity?.length !== 0 &&
				domain.samlConfig?.samlIdp?.length !== 0
			);
		case GOOGLE_AUTH:
			return (
				domain.googleAuthConfig?.clientId?.length !== 0 &&
				domain.googleAuthConfig?.clientSecret?.length !== 0
			);
		case LDAP:
			return (
				domain.ldapConfig?.serverUrl?.length !== 0 &&
				domain.ldapConfig?.bindDn?.length !== 0 &&
				domain.ldapConfig?.bindPassword?.length !== 0 &&
				domain.ldapConfig?.userBaseDn?.length !== 0 &&
				domain.ldapConfig?.userFilter?.length !== 0
			);
		default:
			return false;
	}
};
