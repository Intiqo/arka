package database

import (
	"errors"
	"github.com/adwitiyaio/arka/constants"
	"github.com/jackc/pgconn"
	"gorm.io/gorm"
	"os"
	"testing"

	"github.com/adwitiyaio/arka/config"
	"github.com/adwitiyaio/arka/dependency"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type GormManagerSuite struct {
	suite.Suite

	gdm *gormDatabaseManager
	db  *gorm.DB
}

func (ts *GormManagerSuite) SetupSuite() {
	dm := dependency.GetManager()
	config.Bootstrap(config.ProviderEnvironment, "../test.env")
	ts.gdm = &gormDatabaseManager{
		cm: dm.Get(config.DependencyConfigManager).(config.Manager),
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
		const msg = "You cannot delete a [book] that is associated"
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
