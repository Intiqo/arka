package cloud

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsCfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"

	"github.com/adwitiyaio/arka/config"
	"github.com/adwitiyaio/arka/dependency"
	"github.com/adwitiyaio/arka/logger"
)

const regionKey = "AWS_REGION"
const accessKeyIdKey = "AWS_ACCESS_KEY_ID"
const secretAccessKey = "AWS_SECRET_ACCESS_KEY"

type awsManager struct {
	cm config.Manager

	config aws.Config
	region string
}

func (am *awsManager) GetConfig() aws.Config {
	return am.config
}

func (am *awsManager) GetRegion() string {
	return am.region
}

func (am *awsManager) initialize() {
	dm := dependency.GetManager()
	cm := dm.Get(config.DependencyConfigManager).(config.Manager)

	var err error
	am.config, err = awsCfg.LoadDefaultConfig(context.TODO())
	if err != nil {
		logger.Log.Error().Err(err).Msg("failed to load default aws config. Trying with environemt variables")
		am.region = cm.GetValueForKey(regionKey)
		accessKeyId := cm.GetValueForKey(accessKeyIdKey)
		secretAccessKey := cm.GetValueForKey(secretAccessKey)
		if am.region == "" || accessKeyId == "" || secretAccessKey == "" {
			logger.Log.Fatal().Msg("Missing AWS credentials. Please set AWS_REGION, AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY environment variables")
		}
		am.config = *&aws.Config{
			Region:      am.region,
			Credentials: credentials.NewStaticCredentialsProvider(accessKeyId, secretAccessKey, ""),
		}
	}
}
