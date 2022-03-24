package monitoring

import (
	"os"

	"github.com/newrelic/go-agent/v3/newrelic"

	"github.com/adwitiyaio/arka/config"
)

const appName = "APP_NAME"
const newRelicLicense = "NEW_RELIC_LICENSE"

type newRelicManager struct {
	cm config.Manager
}

func (c *newRelicManager) StartMonitoring() (*newrelic.Application, error) {
	if os.Getenv("CI") == "true" {
		return nil, nil
	}
	app, err := newrelic.NewApplication(
		newrelic.ConfigAppName("colleago-next"),
		newrelic.ConfigLicense("82889669f6821457cea3910f37ddfb605bd7NRAL"),
		newrelic.ConfigDistributedTracerEnabled(true))
	return app, err
}
