package secrets

import (
	"github.com/adwitiyaio/arka/config"
)

type environmentSecretsManager struct {
	cm config.Manager
}

func (esm *environmentSecretsManager) initialize() {}

func (esm *environmentSecretsManager) GetValueForKey(key string) string {
	return esm.cm.GetValueForKey(key)
}
