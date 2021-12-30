package template

import (
	"errors"
	"github.com/adwitiyaio/arka/dependency"
	"github.com/adwitiyaio/arka/logger"
)

const DependencyTemplateManager = "template_manager"

// ProviderSystem is a template provider from go's default template package
const ProviderSystem = "system"

// Manager ... A template manager to create text / html templates
type Manager interface {
	// CreateTemplate ... Create html / text template from the data
	CreateTemplate(ctx string, t string, data interface{}, html bool) (string, error)
}

func Bootstrap(provider string) {
	dm := dependency.GetManager()
	var tm interface{}
	switch provider {
	case ProviderSystem:
		tm = &systemTemplateManager{}
	default:
		err := errors.New("template provider unknown")
		logger.Log.Fatal().Err(err).Str("provider", provider)
	}
	dm.Register(DependencyTemplateManager, tm)
}
