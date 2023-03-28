package util

import (
	"os"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/adwitiyaio/arka/config"
	"github.com/adwitiyaio/arka/dependency"
	"github.com/adwitiyaio/arka/secrets"
)

type UrlManagerTestSuite struct {
	suite.Suite

	m UrlManager
}

func TestUrlManager(t *testing.T) {
	suite.Run(t, new(UrlManagerTestSuite))
}

func (ts *UrlManagerTestSuite) SetupSuite() {
	config.Bootstrap(config.ProviderEnvironment, "../test.env")
	secrets.Bootstrap(secrets.ProviderEnvironment, "")
	err := os.Setenv("CI", "true")
	require.NoError(ts.T(), err)
}

func (ts *UrlManagerTestSuite) Test_multiUrlShortener_Shorten() {
	ts.Run("success - shorten url - kutt", func() {
		BootstrapUrlManager(UrlProviderKutt)
		ts.m = dependency.GetManager().Get(DependencyUrlManager).(UrlManager)
		res, err := ts.m.Shorten(gofakeit.URL())
		assert.NoError(ts.T(), err)
		assert.NotNil(ts.T(), res)
	})
	ts.Run("success - shorten url - smallr links", func() {
		BootstrapUrlManager(UrlProviderSmallrLinks)
		ts.m = dependency.GetManager().Get(DependencyUrlManager).(UrlManager)
		res, err := ts.m.Shorten(gofakeit.URL())
		assert.NoError(ts.T(), err)
		assert.NotNil(ts.T(), res)
	})
}

func (ts *UrlManagerTestSuite) Test_multiUrlShortener_CreateDeepLink() {
	ts.Run("success", func() {
		BootstrapUrlManager(UrlProviderSmallrLinks)
		ts.m = dependency.GetManager().Get(DependencyUrlManager).(UrlManager)
		res, err := ts.m.CreateDeepLink(gofakeit.URL())
		assert.NoError(ts.T(), err)
		assert.NotNil(ts.T(), res)
	})
}
