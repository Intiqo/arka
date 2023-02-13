package secrets

import (
	"errors"

	"github.com/adwitiyaio/arka/cloud"
	"github.com/adwitiyaio/arka/config"
	"github.com/adwitiyaio/arka/dependency"
	"github.com/adwitiyaio/arka/logger"
)

const DependencySecretsManager = "secrets_manager"

const ProviderEnvironment = "ENVIRONMENT"
const ProviderAwsSecrets = "AWS_SECRETS"

const AwsSecretNameKey = "AWS_SECRET_NAME"

type Manager interface {
	// GetValueForKey ... Gets the value for a secret key
	GetValueForKey(key string) string
}

// Bootstrap ... Bootstraps the secrets manager
// Currently supported providers are:
// - ENVIRONMENT
// - AWS_SECRETS
// Path corresponds to the name of the secret for AWS_SECRETS. For ENVIRONMENT, path is ignored.
func Bootstrap(provider string, path string) {
	dm := dependency.GetManager()
	var m interface{}
	switch provider {
	case ProviderEnvironment:
		m = &environmentSecretsManager{
			cm: dm.Get(config.DependencyConfigManager).(config.Manager),
		}
		m.(*environmentSecretsManager).initialize()

	case ProviderAwsSecrets:
		m = &awsSecretsManager{
			path: path,
			clm:  dm.Get(cloud.DependencyCloudManager).(cloud.Manager),
		}
		m.(*awsSecretsManager).initialize()
	default:
		err := errors.New("secrets provider unknown")
		logger.Log.Fatal().Err(err).Str("provider", provider)
	}
	dm.Register(DependencySecretsManager, m)
}
