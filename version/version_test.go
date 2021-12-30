package version

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type VersionTestSuite struct {
	suite.Suite
}

func TestVersion(t *testing.T) {
	suite.Run(t, new(VersionTestSuite))
}

func (ts *VersionTestSuite) Test_PrintInfo() {
	ts.Run("success - print info", func() {
		PrintInfo()
	})
}
