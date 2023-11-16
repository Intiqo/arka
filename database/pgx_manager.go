package database

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/jackc/pgconn"
	pgxuuid "github.com/jackc/pgx-gofrs-uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/adwitiyaio/arka/constants"
	"github.com/adwitiyaio/arka/exception"
	"github.com/adwitiyaio/arka/logger"
	"github.com/adwitiyaio/arka/secrets"
)

type pgxDatabaseManager struct {
	sm secrets.Manager
	db *pgxpool.Pool
}

func (pdm *pgxDatabaseManager) initialize() {
	pdm.db = pdm.connect()
}

func (pdm *pgxDatabaseManager) GetInstance() *pgxpool.Pool {
	return pdm.db
}

func (pdm pgxDatabaseManager) GetStatus() string {
	err := pdm.db.Ping(context.Background())
	if err != nil {
		logger.Log.Error().Err(err).Msg("failed to ping database, connection may be down")
		return constants.SystemStatusDown
	}
	return constants.SystemStatusUp
}

func (pdm pgxDatabaseManager) TranslateError(err error, ent string) error {
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

func (pdm pgxDatabaseManager) connect() *pgxpool.Pool {
	dbHost := strings.TrimSpace(pdm.sm.GetValueForKey(dbHostKey))
	portEnv := strings.TrimSpace(pdm.sm.GetValueForKey(dbPortKey))
	dbPort, err := strconv.Atoi(portEnv)
	if err != nil {
		logger.Log.Warn().Str("port", portEnv).Msg("failed to parse db port from environment, resorting to default")
		dbPort = 5432
	}
	dbName := strings.TrimSpace(pdm.sm.GetValueForKey(dbDatabaseKey))
	dbUser := strings.TrimSpace(pdm.sm.GetValueForKey(dbUserKey))
	dbPassword := strings.TrimSpace(pdm.sm.GetValueForKey(dbPasswordKey))

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbPort, dbName)

	// Create a database connection pool
	dbconfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatalf("Unable to parse DATABASE_URL: %v\n", err)
	}
	// Register the uuid type
	dbconfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		pgxuuid.Register(conn.TypeMap())
		return nil
	}
	// Create the connection pool
	dbpool, err := pgxpool.NewWithConfig(context.Background(), dbconfig)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v\n", err)
	}

	return dbpool
}
