package database

import (
	"errors"
	"github.com/adwitiyaio/arka/config"
	"github.com/adwitiyaio/arka/logger"
	"gorm.io/gorm"

	"github.com/adwitiyaio/arka/dependency"
)

const DependencyDatabaseManager = "database_manager"

const dbHostKey = "DB_HOST"
const dbPortKey = "DB_PORT"
const dbDatabaseKey = "DB_DATABASE"
const dbUserKey = "DB_USER"
const dbPasswordKey = "DB_PASSWORD"

const ProviderGorm = "GORM"

type Manager interface {
	// GetInstance ... Gets an instance of the database
	GetInstance() *gorm.DB

	// GetStatus ... Returns the current status of the database connection
	GetStatus() string

	// TranslateError ... Translates a database exception to user-friendly exception
	TranslateError(err error, ent string) error
}

// Bootstrap ... Bootstraps the database manager
func Bootstrap(providerOrm string) {
	c := dependency.GetManager()
	var dm interface{}
	switch providerOrm {
	case ProviderGorm:
		dm = &gormDatabaseManager{
			cm: c.Get(config.DependencyConfigManager).(config.Manager),
		}
		dm.(*gormDatabaseManager).initialize()
	default:
		err := errors.New("orm provider not implemented")
		logger.Log.Fatal().Err(err).Str("provider", providerOrm)
	}
	c.Register(DependencyDatabaseManager, dm)
}
