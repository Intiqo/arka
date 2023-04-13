package database

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgconn"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gLogger "gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"

	"github.com/adwitiyaio/arka/constants"
	"github.com/adwitiyaio/arka/exception"
	"github.com/adwitiyaio/arka/logger"
	"github.com/adwitiyaio/arka/secrets"
)

const ErrorEntityAssociated = `Not deleted. %s is associated with other records.`

type gormDatabaseManager struct {
	sm secrets.Manager
	db *gorm.DB
}

func (gdm *gormDatabaseManager) initialize() {
	hosts := gdm.sm.GetValueForKey(dbHostsKey)
	if hosts != "" {
		gdm.db = gdm.connectMultiple(hosts)
		return
	}
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
	host := strings.TrimSpace(gdm.sm.GetValueForKey(dbHostKey))
	portEnv := strings.TrimSpace(gdm.sm.GetValueForKey(dbPortKey))
	port, err := strconv.Atoi(portEnv)
	if err != nil {
		logger.Log.Warn().Str("port", portEnv).Msg("failed to parse db port from environment, resorting to default")
		port = 5432
	}
	database := strings.TrimSpace(gdm.sm.GetValueForKey(dbDatabaseKey))
	user := strings.TrimSpace(gdm.sm.GetValueForKey(dbUserKey))
	password := strings.TrimSpace(gdm.sm.GetValueForKey(dbPasswordKey))

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", host, user, password, database, port)

	gormLogger := gLogger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		gLogger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  gLogger.Silent,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      true,
			Colorful:                  false,
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})

	if err != nil {
		logger.Log.Panic().Err(err).Stack().Msg("unable to connect to database")
	}

	return db
}

func (gdm gormDatabaseManager) connectMultiple(hostsStr string) *gorm.DB {
	dbHosts := strings.Split(hostsStr, ",")
	if len(dbHosts) < 1 {
		logger.Log.Panic().Stack().Msg("no hosts provided")
	}

	firstSourceDsn := ""
	sources := make([]gorm.Dialector, 0)
	replicas := make([]gorm.Dialector, 0)
	for _, dbHost := range dbHosts {
		dbHost = strings.TrimSpace(dbHost)

		replParam := strings.Split(dbHost, "?")[1]
		repl := strings.Split(replParam, "=")[1]
		isReplica, err := strconv.ParseBool(repl)
		if err != nil {
			isReplica = false
		}

		database := strings.Split(strings.Split(dbHost, "?")[0], "/")[1]

		host := strings.Split(strings.Split(strings.Split(strings.Split(dbHost, "?")[0], "/")[0], "@")[1], ":")[0]
		portEnv := strings.Split(strings.Split(strings.Split(strings.Split(dbHost, "?")[0], "/")[0], "@")[1], ":")[1]
		port, err := strconv.Atoi(portEnv)
		if err != nil {
			logger.Log.Warn().Str("port", portEnv).Msg("failed to parse db port from environment, resorting to default")
			port = 5432
		}

		user := strings.Split(strings.Split(strings.Split(strings.Split(dbHost, "?")[0], "/")[0], "@")[0], ":")[0]
		password := strings.Split(strings.Split(strings.Split(strings.Split(dbHost, "?")[0], "/")[0], "@")[0], ":")[1]

		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", host, user, password, database, port)

		if isReplica {
			replicas = append(replicas, postgres.Open(dsn))
		} else {
			sources = append(sources, postgres.Open(dsn))
			if firstSourceDsn == "" {
				firstSourceDsn = dsn
			}
		}
	}

	db, err := gorm.Open(postgres.Open(firstSourceDsn), &gorm.Config{})

	if err != nil {
		logger.Log.Panic().Err(err).Stack().Msg("unable to connect to database")
	}

	db.Use(dbresolver.Register(dbresolver.Config{
		Sources:  sources,
		Replicas: replicas,
	}))

	return db
}
