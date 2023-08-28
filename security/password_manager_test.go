package security

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/adwitiyaio/arka/dependency"
)

type PasswordManagerTestSuite struct {
	suite.Suite

	m PasswordManager
}

func (ts *PasswordManagerTestSuite) SetupSuite() {
	BootstrapPasswordManager()
	ts.m = dependency.GetManager().Get(DependencyPasswordManager).(PasswordManager)
}

func TestPasswordManager(t *testing.T) {
	suite.Run(t, new(PasswordManagerTestSuite))
}

func (ts *PasswordManagerTestSuite) Test_simplePasswordManager_HashPassword() {
	ts.Run(
		"success", func() {
			res := ts.m.HashPassword("test")
			assert.NotNil(ts.T(), res)
		},
	)
}

func (ts *PasswordManagerTestSuite) Test_simplePasswordManager_VerifyPassword() {
	ts.Run(
		"success", func() {
			pass := "test"
			hash := ts.m.HashPassword(pass)
			require.NotNil(ts.T(), hash)

			err := ts.m.VerifyPassword(pass, hash)
			require.NoError(ts.T(), err)
		},
	)
}

func (ts *PasswordManagerTestSuite) Test_simplePasswordManager_CreateSha1Password() {
	ts.Run(
		"success", func() {
			res := ts.m.CreateSha1Hash("test")
			assert.NotNil(ts.T(), res)
		},
	)
}
