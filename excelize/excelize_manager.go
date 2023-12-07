package excelize

import (
	"github.com/adwitiyaio/arka/dependency"
	"github.com/xuri/excelize/v2"
)

const DependencyExcelizeManager = "excelize_manager"

type Manager interface {
	// NewFile ... Creates a new excelize file
	NewFile(opts excelize.Options) *excelize.File
}

type Excelize struct {
}

// Bootstrap ... Bootstraps the excelize manager
func Bootstrap() {
	d := dependency.GetManager()
	e := &Excelize{}
	d.Register(DependencyExcelizeManager, e)
}

func (c *Excelize) NewFile(opts excelize.Options) *excelize.File {
	return excelize.NewFile(opts)
}
