package database

import (
	"errors"
	"os"
	"testing"

	"github.com/jackc/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"

	"github.com/adwitiyaio/arka/config"
	"github.com/adwitiyaio/arka/constants"
	"github.com/adwitiyaio/arka/dependency"
	"github.com/adwitiyaio/arka/secrets"
)

type GormManagerSuite struct {
	suite.Suite

	gdm *gormDatabaseManager
	db  *gorm.DB
}

func (ts *GormManagerSuite) SetupSuite() {
	dm := dependency.GetManager()
	config.Bootstrap(config.ProviderEnvironment, "../test.env")
	secrets.Bootstrap(secrets.ProviderEnvironment, "")
	ts.gdm = &gormDatabaseManager{
		sm: dm.Get(secrets.DependencySecretsManager).(secrets.Manager),
	}
	ts.gdm.db = ts.gdm.connect()

	Bootstrap(ProviderGorm)
}

func TestGormManager(t *testing.T) {
	suite.Run(t, new(GormManagerSuite))
}

func (ts GormManagerSuite) Test_Connect() {
	ts.Run("invalid port from config", func() {
		err := os.Setenv(dbPortKey, "unknown")
		require.NoError(ts.T(), err)
		db := ts.gdm.connect()
		assert.NotNil(ts.T(), db)
	})
}

func (ts GormManagerSuite) Test_GetInstance() {
	ts.Run("success", func() {
		db := ts.gdm.GetInstance()
		assert.NotNil(ts.T(), db)
	})
}

func (ts GormManagerSuite) Test_GetStatus() {
	ts.Run("success - UP", func() {
		status := ts.gdm.GetStatus()
		assert.Equal(ts.T(), constants.SystemStatusUp, status)
	})
}

func (ts GormManagerSuite) Test_TranslateError() {
	ts.Run("success - entity associated", func() {
		pgErr := pgconn.PgError{Code: "23503"}
		err := ts.gdm.TranslateError(&pgErr, "book")
		const msg = "Not deleted. [Book] is associated with other records."
		assert.Equal(ts.T(), msg, err.Error())
	})
	ts.Run("success - unknown code", func() {
		pgErr := pgconn.PgError{Code: "12345", Message: "Unknown error"}
		err := ts.gdm.TranslateError(&pgErr, "book")
		assert.Equal(ts.T(), pgErr.Message, err.Error())
	})
	ts.Run("success - non pg error", func() {
		genericErr := errors.New("generic error")
		err := ts.gdm.TranslateError(genericErr, "book")
		assert.Equal(ts.T(), genericErr.Error(), err.Error())
	})
}
