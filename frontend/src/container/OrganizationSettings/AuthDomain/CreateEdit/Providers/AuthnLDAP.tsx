import { Style } from '@signozhq/design-tokens';
import { CircleHelp } from '@signozhq/icons';
import { Checkbox, Input } from '@signozhq/ui';
import { Form, InputNumber, Tooltip } from 'antd';

import './Providers.styles.scss';

function ConfigureLDAPAuthnProvider({
	isCreate,
}: {
	isCreate: boolean;
}): JSX.Element {
	const form = Form.useFormInstance();

	return (
		<div className="authn-provider">
			<section className="authn-provider__header">
				<h3 className="authn-provider__title">Configure LDAP Authentication</h3>
				<p className="authn-provider__description">
					Authenticate users against your LDAP directory. Users will sign in with
					their email address and LDAP password using the standard SigNoz login
					form.
				</p>
			</section>

			<div className="authn-provider__columns">
				<div className="authn-provider__left">
					<div className="authn-provider__field-group">
						<label className="authn-provider__label" htmlFor="ldap-domain">
							Domain
							<Tooltip title="The email domain for users who should use LDAP (e.g., example.com for users with @example.com emails)">
								<CircleHelp size={14} color={Style.L3_FOREGROUND} cursor="help" />
							</Tooltip>
						</label>
						<Form.Item
							name="name"
							className="authn-provider__form-item"
							rules={[
								{ required: true, message: 'Domain is required', whitespace: true },
							]}
						>
							<Input id="ldap-domain" disabled={!isCreate} />
						</Form.Item>
					</div>

					<div className="authn-provider__field-group">
						<label className="authn-provider__label" htmlFor="ldap-server-url">
							LDAP Server URL
							<Tooltip title="Protocol and hostname (e.g., ldap://ldap.example.com or ldaps://ldap.example.com)">
								<CircleHelp size={14} color={Style.L3_FOREGROUND} cursor="help" />
							</Tooltip>
						</label>
						<Form.Item
							name={['ldapConfig', 'serverUrl']}
							className="authn-provider__form-item"
							rules={[
								{ required: true, message: 'Server URL is required', whitespace: true },
							]}
						>
							<Input id="ldap-server-url" placeholder="ldap://ldap.example.com" />
						</Form.Item>
					</div>

					<div className="authn-provider__field-group">
						<label className="authn-provider__label" htmlFor="ldap-server-port">
							Port
							<Tooltip title="Default: 389 for LDAP, 636 for LDAPS">
								<CircleHelp size={14} color={Style.L3_FOREGROUND} cursor="help" />
							</Tooltip>
						</label>
						<Form.Item
							name={['ldapConfig', 'serverPort']}
							className="authn-provider__form-item"
							initialValue={389}
							rules={[{ required: true, message: 'Port is required' }]}
						>
							<InputNumber
								id="ldap-server-port"
								min={1}
								max={65535}
								style={{ width: '100%' }}
							/>
						</Form.Item>
					</div>

					<div className="authn-provider__field-group">
						<label className="authn-provider__label" htmlFor="ldap-bind-dn">
							Bind DN
							<Tooltip title="Distinguished Name of the service account (e.g., cn=admin,dc=example,dc=com)">
								<CircleHelp size={14} color={Style.L3_FOREGROUND} cursor="help" />
							</Tooltip>
						</label>
						<Form.Item
							name={['ldapConfig', 'bindDn']}
							className="authn-provider__form-item"
							rules={[{ required: true, message: 'Bind DN is required', whitespace: true }]}
						>
							<Input
								id="ldap-bind-dn"
								placeholder="cn=admin,dc=example,dc=com"
							/>
						</Form.Item>
					</div>

					<div className="authn-provider__field-group">
						<label className="authn-provider__label" htmlFor="ldap-bind-password">
							Bind Password
						</label>
						<Form.Item
							name={['ldapConfig', 'bindPassword']}
							className="authn-provider__form-item"
							rules={[
								{ required: true, message: 'Bind password is required' },
							]}
						>
							<Input.Password
								id="ldap-bind-password"
								placeholder="Service account password"
							/>
						</Form.Item>
					</div>

					<div className="authn-provider__field-group">
						<label className="authn-provider__label" htmlFor="ldap-user-base-dn">
							User Base DN
							<Tooltip title="Base DN for user searches (e.g., ou=users,dc=example,dc=com)">
								<CircleHelp size={14} color={Style.L3_FOREGROUND} cursor="help" />
							</Tooltip>
						</label>
						<Form.Item
							name={['ldapConfig', 'userBaseDn']}
							className="authn-provider__form-item"
							rules={[
								{ required: true, message: 'User Base DN is required', whitespace: true },
							]}
						>
							<Input
								id="ldap-user-base-dn"
								placeholder="ou=users,dc=example,dc=com"
							/>
						</Form.Item>
					</div>

					<div className="authn-provider__field-group">
						<label className="authn-provider__label" htmlFor="ldap-user-filter">
							User Filter
							<Tooltip title="LDAP filter for user search. Use %s as placeholder replaced with the user's identifier (e.g., (mail=%s) or (sAMAccountName=%s))">
								<CircleHelp size={14} color={Style.L3_FOREGROUND} cursor="help" />
							</Tooltip>
						</label>
						<Form.Item
							name={['ldapConfig', 'userFilter']}
							className="authn-provider__form-item"
							rules={[
								{ required: true, message: 'User filter is required', whitespace: true },
							]}
						>
							<Input id="ldap-user-filter" placeholder="(mail=%s)" />
						</Form.Item>
					</div>
				</div>

				<div className="authn-provider__right">
					<div className="authn-provider__field-group">
						<label
							className="authn-provider__label"
							htmlFor="ldap-email-attribute"
						>
							Email Attribute
							<Tooltip title="LDAP attribute containing the user's email address (default: mail)">
								<CircleHelp size={14} color={Style.L3_FOREGROUND} cursor="help" />
							</Tooltip>
						</label>
						<Form.Item
							name={['ldapConfig', 'emailAttribute']}
							className="authn-provider__form-item"
							initialValue="mail"
						>
							<Input id="ldap-email-attribute" placeholder="mail" />
						</Form.Item>
					</div>

					<div className="authn-provider__field-group">
						<label
							className="authn-provider__label"
							htmlFor="ldap-display-name-attribute"
						>
							Display Name Attribute
							<Tooltip title="LDAP attribute for the user's display name (default: cn)">
								<CircleHelp size={14} color={Style.L3_FOREGROUND} cursor="help" />
							</Tooltip>
						</label>
						<Form.Item
							name={['ldapConfig', 'displayNameAttribute']}
							className="authn-provider__form-item"
							initialValue="cn"
						>
							<Input id="ldap-display-name-attribute" placeholder="cn" />
						</Form.Item>
					</div>

					<div className="authn-provider__checkbox-row">
						<Form.Item
							name={['ldapConfig', 'useTls']}
							valuePropName="value"
							noStyle
							initialValue={false}
						>
							<Checkbox
								id="ldap-use-tls"
								onChange={(checked: boolean): void => {
									form.setFieldValue(['ldapConfig', 'useTls'], checked);
								}}
							>
								Use TLS (ldaps://)
							</Checkbox>
						</Form.Item>
					</div>

					<div className="authn-provider__checkbox-row">
						<Form.Item
							name={['ldapConfig', 'useStartTls']}
							valuePropName="value"
							noStyle
							initialValue={false}
						>
							<Checkbox
								id="ldap-use-start-tls"
								onChange={(checked: boolean): void => {
									form.setFieldValue(['ldapConfig', 'useStartTls'], checked);
								}}
							>
								Use StartTLS
							</Checkbox>
						</Form.Item>
					</div>

					<div className="authn-provider__checkbox-row">
						<Form.Item
							name={['ldapConfig', 'skipTlsVerify']}
							valuePropName="value"
							noStyle
							initialValue={false}
						>
							<Checkbox
								id="ldap-skip-tls-verify"
								onChange={(checked: boolean): void => {
									form.setFieldValue(['ldapConfig', 'skipTlsVerify'], checked);
								}}
							>
								Skip TLS Verification
								<Tooltip title="Insecure — disable only for testing">
									<CircleHelp
										size={14}
										color={Style.L3_FOREGROUND}
										cursor="help"
										style={{ marginLeft: 4 }}
									/>
								</Tooltip>
							</Checkbox>
						</Form.Item>
					</div>

					<div className="authn-provider__field-group" style={{ marginTop: 16 }}>
						<label
							className="authn-provider__label"
							htmlFor="ldap-group-base-dn"
						>
							Group Base DN (Optional)
							<Tooltip title="Base DN for group searches (e.g., ou=groups,dc=example,dc=com)">
								<CircleHelp size={14} color={Style.L3_FOREGROUND} cursor="help" />
							</Tooltip>
						</label>
						<Form.Item
							name={['ldapConfig', 'groupBaseDn']}
							className="authn-provider__form-item"
						>
							<Input
								id="ldap-group-base-dn"
								placeholder="ou=groups,dc=example,dc=com"
							/>
						</Form.Item>
					</div>

					<div className="authn-provider__field-group">
						<label
							className="authn-provider__label"
							htmlFor="ldap-group-filter"
						>
							Group Filter (Optional)
							<Tooltip title="LDAP filter for group search. Use %s as placeholder for user DN (e.g., (member=%s))">
								<CircleHelp size={14} color={Style.L3_FOREGROUND} cursor="help" />
							</Tooltip>
						</label>
						<Form.Item
							name={['ldapConfig', 'groupFilter']}
							className="authn-provider__form-item"
						>
							<Input id="ldap-group-filter" placeholder="(member=%s)" />
						</Form.Item>
					</div>
				</div>
			</div>
		</div>
	);
}

export default ConfigureLDAPAuthnProvider;
