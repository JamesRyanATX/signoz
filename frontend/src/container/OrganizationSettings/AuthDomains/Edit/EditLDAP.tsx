import { InfoCircleFilled } from '@ant-design/icons';
import {
	Card,
	Checkbox,
	Collapse,
	Form,
	Input,
	InputNumber,
	Space,
	Typography,
} from 'antd';

const { Panel } = Collapse;

function EditLDAP(): JSX.Element {
	return (
		<>
			<Typography.Title level={5}>Server Configuration</Typography.Title>

			<Form.Item
				label="LDAP Server URL"
				name={['ldapConfig', 'serverUrl']}
				rules={[
					{
						required: true,
						message: 'Please input LDAP server URL (e.g., ldap://ldap.example.com)',
					},
				]}
				tooltip="Protocol and hostname (e.g., ldap://ldap.example.com or ldaps://ldap.example.com)"
			>
				<Input placeholder="ldap://ldap.example.com" />
			</Form.Item>

			<Form.Item
				label="Port"
				name={['ldapConfig', 'serverPort']}
				rules={[{ required: true, message: 'Please input port number' }]}
				tooltip="Default: 389 for LDAP, 636 for LDAPS"
				initialValue={389}
			>
				<InputNumber min={1} max={65535} style={{ width: '100%' }} />
			</Form.Item>

			<Typography.Title level={5}>Bind Configuration</Typography.Title>

			<Form.Item
				label="Bind DN"
				name={['ldapConfig', 'bindDn']}
				rules={[{ required: true, message: 'Please input bind DN' }]}
				tooltip="Distinguished Name of the service account (e.g., cn=admin,dc=example,dc=com)"
			>
				<Input placeholder="cn=admin,dc=example,dc=com" />
			</Form.Item>

			<Form.Item
				label="Bind Password"
				name={['ldapConfig', 'bindPassword']}
				rules={[{ required: true, message: 'Please input bind password' }]}
			>
				<Input.Password placeholder="Service account password" />
			</Form.Item>

			<Typography.Title level={5}>User Search Configuration</Typography.Title>

			<Form.Item
				label="User Base DN"
				name={['ldapConfig', 'userBaseDn']}
				rules={[{ required: true, message: 'Please input user base DN' }]}
				tooltip="Base DN for user searches (e.g., ou=users,dc=example,dc=com)"
			>
				<Input placeholder="ou=users,dc=example,dc=com" />
			</Form.Item>

			<Form.Item
				label="User Filter"
				name={['ldapConfig', 'userFilter']}
				rules={[{ required: true, message: 'Please input user filter' }]}
				tooltip="LDAP filter for user search. Use %s as placeholder which will be replaced with the user's email address (case-insensitive). Examples: (mail=%s) for email-based search, or (sAMAccountName=%s) if email prefix matches username"
			>
				<Input placeholder="(mail=%s)" />
			</Form.Item>

			<Form.Item
				label="Username Attribute"
				name={['ldapConfig', 'usernameAttribute']}
				initialValue="uid"
				tooltip="LDAP attribute for username (default: uid)"
			>
				<Input placeholder="uid" />
			</Form.Item>

			<Form.Item
				label="Email Attribute"
				name={['ldapConfig', 'emailAttribute']}
				initialValue="mail"
				tooltip="LDAP attribute for email (default: mail)"
			>
				<Input placeholder="mail" />
			</Form.Item>

			<Form.Item
				label="Display Name Attribute"
				name={['ldapConfig', 'displayNameAttribute']}
				initialValue="cn"
				tooltip="LDAP attribute for display name (default: cn)"
			>
				<Input placeholder="cn" />
			</Form.Item>

			<Collapse ghost>
				<Panel header="Advanced Settings (Optional)" key="1">
					<Typography.Title level={5}>Security Settings</Typography.Title>

					<Form.Item
						name={['ldapConfig', 'useTls']}
						valuePropName="checked"
						initialValue={false}
					>
						<Checkbox>Use TLS (ldaps://)</Checkbox>
					</Form.Item>

					<Form.Item
						name={['ldapConfig', 'useStartTls']}
						valuePropName="checked"
						initialValue={false}
					>
						<Checkbox>Use StartTLS</Checkbox>
					</Form.Item>

					<Form.Item
						name={['ldapConfig', 'skipTlsVerify']}
						valuePropName="checked"
						initialValue={false}
					>
						<Checkbox>Skip TLS Verification (insecure, for testing only)</Checkbox>
					</Form.Item>

					<Typography.Title level={5}>Group Search (Optional)</Typography.Title>

					<Form.Item
						label="Group Base DN"
						name={['ldapConfig', 'groupBaseDn']}
						tooltip="Base DN for group searches (e.g., ou=groups,dc=example,dc=com)"
					>
						<Input placeholder="ou=groups,dc=example,dc=com" />
					</Form.Item>

					<Form.Item
						label="Group Filter"
						name={['ldapConfig', 'groupFilter']}
						tooltip="LDAP filter for group search. Use %s as placeholder for user DN (e.g., (member=%s))"
					>
						<Input placeholder="(member=%s)" />
					</Form.Item>

					<Form.Item
						label="Group Member Attribute"
						name={['ldapConfig', 'groupMemberAttr']}
						initialValue="member"
					>
						<Input placeholder="member" />
					</Form.Item>

					<Typography.Title level={5}>Timeouts (Optional)</Typography.Title>

					<Form.Item
						label="Search Timeout (seconds)"
						name={['ldapConfig', 'searchTimeout']}
						initialValue={10}
					>
						<InputNumber min={1} max={60} style={{ width: '100%' }} />
					</Form.Item>

					<Form.Item
						label="Connection Timeout (seconds)"
						name={['ldapConfig', 'connTimeout']}
						initialValue={10}
					>
						<InputNumber min={1} max={60} style={{ width: '100%' }} />
					</Form.Item>
				</Panel>
			</Collapse>

			<Card style={{ marginTop: '1rem' }}>
				<Space>
					<InfoCircleFilled />
					<Typography>
						LDAP authentication will be enabled for users matching this domain. Users
						will log in with their email address and password using the standard
						SigNoz login form. The User Filter determines how the email is searched in
						LDAP.
					</Typography>
				</Space>
			</Card>
		</>
	);
}

export default EditLDAP;
