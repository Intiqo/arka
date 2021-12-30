package dependency

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type testService interface {
	Greet(message string)
}

type testServiceImpl struct {
}

func (t testServiceImpl) Greet(message string) {
	fmt.Println(message)
}

type LocalDependencyManagerTestSuite struct {
	suite.Suite

	cc Manager
}

func (ts *LocalDependencyManagerTestSuite) SetupSuite() {
	ts.cc = GetManager()
}

func TestLocalDependencyManager(t *testing.T) {
	suite.Run(t, new(LocalDependencyManagerTestSuite))
}

func (ts LocalDependencyManagerTestSuite) Test_localDependencyManager_RegisterAndGet() {
	ts.Run("success", func() {
		ts.cc.Register("test_service", &testServiceImpl{})

		as := ts.cc.Get("test_service").(testService)
		require.NotNil(ts.T(), as)

		testService := &testServiceImpl{}
		assert.IsType(ts.T(), testService, as)
	})

}
