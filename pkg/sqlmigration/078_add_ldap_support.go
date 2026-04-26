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
	// No schema changes needed. LDAP configuration is stored in auth_domain.data as JSON.
	return nil
}

func (migration *addLdapSupport) Down(context.Context, *bun.DB) error {
	return nil
}
