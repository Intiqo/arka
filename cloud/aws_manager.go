package cloud

import (
	"context"
	"log"

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
	// Try to load credentials from the specified aws profile
	awsProfile := cm.GetValueForKey("AWS_PROFILE")
	logger.Log.Info().Msg("Loading AWS credentials from profile")
	am.config, err = awsCfg.LoadDefaultConfig(context.TODO(), awsCfg.WithSharedConfigProfile(awsProfile))
	if err != nil {
		// Fallback to loading credentials from environment variables
		logger.Log.Info().Msg("Loading AWS credentials from environment variables")
		am.region = cm.GetValueForKey(regionKey)
		accessKeyId := cm.GetValueForKey(accessKeyIdKey)
		secretAccessKey := cm.GetValueForKey(secretAccessKey)
		if am.region == "" || accessKeyId == "" || secretAccessKey == "" {
			log.Fatalf("Missing AWS credentials. Please set AWS_REGION, AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY environment variables")
		}
		am.config = aws.Config{
			Region:      am.region,
			Credentials: credentials.NewStaticCredentialsProvider(accessKeyId, secretAccessKey, ""),
		}
	}
}
