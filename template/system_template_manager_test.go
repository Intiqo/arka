package template

import (
	"os"
	"testing"

	"github.com/adwitiyaio/arka/config"
	"github.com/adwitiyaio/arka/dependency"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type SystemTemplateManagerTestSuite struct {
	suite.Suite
	m Manager
}

func TestSystemTemplateManager(t *testing.T) {
	suite.Run(t, new(SystemTemplateManagerTestSuite))
}

func (ts *SystemTemplateManagerTestSuite) SetupSuite() {
	config.Bootstrap(config.ProviderEnvironment, "../test.env")
	err := os.Setenv("CI", "true")
	require.NoError(ts.T(), err)
	Bootstrap(ProviderSystem)
	ts.m = dependency.GetManager().Get(DependencyTemplateManager).(Manager)
}

func (ts *SystemTemplateManagerTestSuite) Test_systemTemplateManager_CreateTemplateForData() {
	ts.Run("error", func() {
		tmp := "Hello {{.Name}}, {{.Age}}"
		obj := struct {
			Name string
		}{
			Name: "Jon Snow",
		}
		_, err := ts.m.CreateTemplate("test - name", tmp, obj, false)
		assert.Error(ts.T(), err)
	})

	ts.Run("success - text", func() {
		tmp := "Hello {{.Name}}"
		obj := struct {
			Name string
		}{
			Name: "Jon Snow",
		}
		result, err := ts.m.CreateTemplate("test - name", tmp, obj, false)
		require.NoError(ts.T(), err)
		require.NotNil(ts.T(), result)

		assert.Equal(ts.T(), "Hello Jon Snow", result)
	})

	ts.Run("success - html", func() {
		tmp := "<body>Hello {{.Name}}</body>"
		obj := struct {
			Name string
		}{
			Name: "Jon Snow",
		}
		result, err := ts.m.CreateTemplate("test - name", tmp, obj, true)
		require.NoError(ts.T(), err)
		require.NotNil(ts.T(), result)

		assert.Equal(ts.T(), "<body>Hello Jon Snow</body>", result)
	})
}
