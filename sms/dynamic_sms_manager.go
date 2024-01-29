package sms

import (
	"github.com/adwitiyaio/arka/cloud"
	"github.com/adwitiyaio/arka/dependency"
	"github.com/adwitiyaio/arka/secrets"
)

type dynamicSmsManager struct {
	msmsm *multiSmsManager
	ssm   *snsManager
	bsmsm *burstSmsManager
}

func (m *dynamicSmsManager) initialize() {
	dm := dependency.GetManager()

	m.msmsm = &multiSmsManager{
		sm: dm.Get(secrets.DependencySecretsManager).(secrets.Manager),
	}
	m.msmsm.initialize()

	m.ssm = &snsManager{
		clm: dm.Get(cloud.DependencyCloudManager).(cloud.Manager),
	}
	m.ssm.initialize()

	m.bsmsm = &burstSmsManager{
		sm: dm.Get(secrets.DependencySecretsManager).(secrets.Manager),
	}
	m.bsmsm.initialize()
}

func (m *dynamicSmsManager) SendSms(options Options) (interface{}, error) {
	switch options.Provider {
	case ProviderMulti:
		return m.msmsm.SendSms(options)
	case ProviderSns:
		return m.ssm.SendSms(options)
	case ProviderBurstSms:
		return m.bsmsm.SendSms(options)
	default:
		return m.msmsm.SendSms(options)
	}
}
