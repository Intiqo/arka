package database

import (
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
	"gorm.io/gorm"

	"github.com/adwitiyaio/arka/dependency"
	"github.com/adwitiyaio/arka/logger"
	"github.com/adwitiyaio/arka/secrets"
)

const DependencyDatabaseManager = "database_manager"

const dbHostKey = "DB_HOST"
const dbPortKey = "DB_PORT"
const dbDatabaseKey = "DB_DATABASE"
const dbUserKey = "DB_USER"
const dbPasswordKey = "DB_PASSWORD"

const dbHostsKey = "DB_HOSTS"

const ProviderGorm = "GORM"
const ProviderPgx = "PGX"

type TDatabase interface {
	*gorm.DB | *pgxpool.Pool
}

type Manager[T TDatabase] interface {
	// GetInstance ... Gets an instance of the database
	GetInstance() T

	// GetStatus ... Returns the current status of the database connection
	GetStatus() string

	// TranslateError ... Translates a database exception to user-friendly exception
	TranslateError(err error, ent string) error
}

// Bootstrap ... Bootstraps the database manager
// If you need to connect to multiple database hosts, check the sample env config for the key `DB_HOSTS`
func Bootstrap(providerOrm string) {
	c := dependency.GetManager()
	var dm interface{}
	switch providerOrm {
	case ProviderGorm:
		dm = &gormDatabaseManager{
			sm: c.Get(secrets.DependencySecretsManager).(secrets.Manager),
		}
		dm.(*gormDatabaseManager).initialize()
	case ProviderPgx:
		dm = &pgxDatabaseManager{
			sm: c.Get(secrets.DependencySecretsManager).(secrets.Manager),
		}
		dm.(*pgxDatabaseManager).initialize()
	default:
		err := errors.New("orm provider not implemented")
		logger.Log.Fatal().Err(err).Str("provider", providerOrm)
	}
	c.Register(DependencyDatabaseManager, dm)
}
