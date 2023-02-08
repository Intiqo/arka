package cloud

import (
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"

	"github.com/adwitiyaio/arka/config"
	"github.com/adwitiyaio/arka/dependency"
	"github.com/adwitiyaio/arka/logger"
)

const DependencyCloudManager = "cloud_manager"

const ProviderAws = "AWS"

type Manager interface {
	GetConfig() aws.Config
	// GetRegion Gets the region for the cloud provider
	GetRegion() string
}

// Bootstrap ... Bootstraps the cloud manager
func Bootstrap(provider string) {
	dm := dependency.GetManager()
	var cm interface{}
	switch provider {
	case ProviderAws:
		cm = &awsManager{
			cm: dm.Get(config.DependencyConfigManager).(config.Manager),
		}
		cm.(*awsManager).initialize()
	default:
		err := errors.New("cloud provider not implemented")
		logger.Log.Fatal().Err(err).Str("provider", provider)
	}
	dm.Register(DependencyCloudManager, cm)
}
