package secrets

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"

	"github.com/adwitiyaio/arka/cloud"
	"github.com/adwitiyaio/arka/logger"
)

type awsSecretsManager struct {
	path string

	clm    cloud.Manager
	client *secretsmanager.Client

	secretsMap map[string]string
}

func (scm *awsSecretsManager) initialize() {
	scm.client = secretsmanager.NewFromConfig(scm.clm.GetConfig())
	valOut, err := scm.client.GetSecretValue(context.Background(), &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(scm.path),
	})
	if err != nil {
		logger.Log.Fatal().Err(err).Stack().Msg("unable to get secret value")
	}
	err = json.Unmarshal([]byte(*valOut.SecretString), &scm.secretsMap)
	if err != nil {
		logger.Log.Fatal().Err(err).Stack().Msg("unable to unmarshal secret value")
	}
}

func (scm *awsSecretsManager) GetValueForKey(key string) string {
	return scm.secretsMap[key]
}
