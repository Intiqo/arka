package database

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/adwitiyaio/arka/exception"

	"github.com/adwitiyaio/arka/constants"

	"github.com/jackc/pgconn"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/adwitiyaio/arka/config"
	"github.com/adwitiyaio/arka/logger"
)

const ErrorEntityAssociated = `Not deleted. %s is associated with other records.`

type gormDatabaseManager struct {
	cm config.Manager
	db *gorm.DB
}

func (gdm *gormDatabaseManager) initialize() {
	gdm.db = gdm.connect()
}

func (gdm *gormDatabaseManager) GetInstance() *gorm.DB {
	return gdm.db
}

func (gdm gormDatabaseManager) GetStatus() string {
	sqlDb, err := gdm.db.DB()
	if err != nil {
		logger.Log.Error().Err(err).Msg("failed to get database instance, connection may be down")
		return constants.SystemStatusUnknown
	}
	err = sqlDb.Ping()
	if err != nil {
		logger.Log.Error().Err(err).Msg("failed to ping database, connection may be down")
		return constants.SystemStatusDown
	}
	return constants.SystemStatusUp
}

func (gdm gormDatabaseManager) TranslateError(err error, ent string) error {
	const entityAssociatedCode = "23503"
	if err, ok := err.(*pgconn.PgError); ok {
		switch err.Code {
		case entityAssociatedCode:
			return exception.CreateAppException(ErrorEntityAssociated, strings.Title(ent))
		default:
			return exception.CreateAppException(err.Message)
		}
	}
	return err
}

func (gdm gormDatabaseManager) connect() *gorm.DB {
	host := strings.TrimSpace(gdm.cm.GetValueForKey(dbHostKey))
	portEnv := strings.TrimSpace(gdm.cm.GetValueForKey(dbPortKey))
	port, err := strconv.Atoi(portEnv)
	if err != nil {
		logger.Log.Warn().Str("port", portEnv).Msg("failed to parse db port from environment, resorting to default")
		port = 5432
	}
	database := strings.TrimSpace(gdm.cm.GetValueForKey(dbDatabaseKey))
	user := strings.TrimSpace(gdm.cm.GetValueForKey(dbUserKey))
	password := strings.TrimSpace(gdm.cm.GetValueForKey(dbPasswordKey))

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", host, user, password, database, port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		logger.Log.Panic().Err(err).Stack().Msg("unable to connect to database")
	}

	return db
}
