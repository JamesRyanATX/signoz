package sqlmigration

import (
	"context"

	"github.com/SigNoz/signoz/pkg/factory"
	"github.com/SigNoz/signoz/pkg/sqlstore"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/migrate"
)

type addLdapSupport struct {
	store sqlstore.SQLStore
}

func NewAddLdapSupportFactory(sqlstore sqlstore.SQLStore) factory.ProviderFactory[SQLMigration, Config] {
	return factory.NewProviderFactory(factory.MustNewName("add_ldap_support"), func(ctx context.Context, ps factory.ProviderSettings, c Config) (SQLMigration, error) {
		return newAddLdapSupport(ctx, ps, c, sqlstore)
	})
}

func newAddLdapSupport(_ context.Context, _ factory.ProviderSettings, _ Config, store sqlstore.SQLStore) (SQLMigration, error) {
	return &addLdapSupport{store: store}, nil
}

func (migration *addLdapSupport) Register(migrations *migrate.Migrations) error {
	if err := migrations.Register(migration.Up, migration.Down); err != nil {
		return err
	}

	return nil
}

func (migration *addLdapSupport) Up(ctx context.Context, db *bun.DB) error {
	// No schema changes are needed for LDAP support
	// The org_domains table already has a `data` JSON column that can store LDAP configuration
	// This migration exists to document the feature addition and ensure proper ordering

	// The LDAP configuration is stored in the org_domains.data JSON field with the following structure:
	// {
	//   "ssoEnabled": true,
	//   "ssoType": "LDAP",
	//   "ldapConfig": {
	//     "serverUrl": "ldap://ldap.example.com",
	//     "serverPort": 389,
	//     "bindDn": "cn=admin,dc=example,dc=com",
	//     "bindPassword": "encrypted_password",
	//     "userBaseDn": "ou=users,dc=example,dc=com",
	//     "userFilter": "(uid=%s)",
	//     "usernameAttribute": "uid",
	//     "emailAttribute": "mail",
	//     "displayNameAttribute": "cn",
	//     "groupBaseDn": "ou=groups,dc=example,dc=com",
	//     "groupFilter": "(member=%s)",
	//     "groupMemberAttr": "member",
	//     "useTls": false,
	//     "useStartTls": true,
	//     "skipTlsVerify": false,
	//     "searchTimeout": 10,
	//     "connTimeout": 10
	//   }
	// }

	return nil
}

func (migration *addLdapSupport) Down(context.Context, *bun.DB) error {
	// No rollback needed as no schema changes were made
	return nil
}
