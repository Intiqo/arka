package config

import (
	"github.com/adwitiyaio/arka/logger"
	"github.com/joho/godotenv"
	"os"
)

type environmentConfigManager struct{}

func (cm *environmentConfigManager) initialize(configPath string) {
	err := godotenv.Load(configPath)
	if err != nil {
		logger.Log.Panic().Err(err).Stack().Msg("unable to initialize configuration")
	}
}

func (cm environmentConfigManager) GetValueForKey(key string) string {
	return os.Getenv(key)
}
