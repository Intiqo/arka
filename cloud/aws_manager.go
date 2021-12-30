package cloud

import (
	"github.com/adwitiyaio/arka/config"
	"github.com/adwitiyaio/arka/dependency"
	"github.com/adwitiyaio/arka/logger"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

const regionKey = "AWS_REGION"
const accessKeyIdKey = "AWS_ACCESS_KEY_ID"
const secretAccessKey = "AWS_SECRET_ACCESS_KEY"

type awsManager struct {
	cm config.Manager

	session *session.Session
	region  string
}

func (am *awsManager) GetSession() *session.Session {
	return am.session
}

func (am *awsManager) GetRegion() string {
	return am.region
}

func (am *awsManager) initialize() {
	dm := dependency.GetManager()
	cm := dm.Get(config.DependencyConfigManager).(config.Manager)

	am.region = cm.GetValueForKey(regionKey)
	accessKeyId := cm.GetValueForKey(accessKeyIdKey)
	secretAccessKey := cm.GetValueForKey(secretAccessKey)

	var err error
	am.session, err = session.NewSession(
		&aws.Config{
			Region: aws.String(am.region),
			Credentials: credentials.NewStaticCredentials(
				accessKeyId, secretAccessKey, "",
			),
		},
	)

	if err != nil {
		logger.Log.Fatal().Err(err).Msg("failed to connect to aws")
	}
}
