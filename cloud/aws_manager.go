package cloud

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsCfg "github.com/aws/aws-sdk-go-v2/config"

	"github.com/adwitiyaio/arka/config"
	"github.com/adwitiyaio/arka/dependency"
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
		awsProfile := cm.GetValueForKey("AWS_PROFILE")
		if awsProfile == "" {
			awsProfile = "default"
		}
		am.config, err = awsCfg.LoadDefaultConfig(context.TODO(), awsCfg.WithSharedConfigProfile(awsProfile))
		if err != nil {
			panic(err)
		}
	}
}
