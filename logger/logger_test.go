package logger

import (
	"github.com/stretchr/testify/require"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

type LoggerTestSuite struct {
	suite.Suite
}

func TestLogger(t *testing.T) {
	suite.Run(t, new(LoggerTestSuite))
}

func (ts *LoggerTestSuite) Test_Bootstrap() {
	ts.Run("success", func() {
		Bootstrap()
	})
	ts.Run("success - production", func() {
		err := os.Setenv("APP_PRODUCTION", "true")
		require.NoError(ts.T(), err)
		Bootstrap()
	})
}
