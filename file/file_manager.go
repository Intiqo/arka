package file

import (
	"github.com/adwitiyaio/arka/cloud"
	"github.com/adwitiyaio/arka/config"
	"io"

	"github.com/adwitiyaio/arka/dependency"
)

const ErrorUnsupportedFileType = "This is an unsupported file type"

const DependencyFileManager = "file_manager"
const ProviderAwsS3 = "aws_s3"

type Manager interface {
	// UploadFile ... Upload file to a storage bucket and returns the corresponding url
	UploadFile(filename string, contentType string, file io.Reader, directory *string) (string, error)

	// GetExtensionAndContentType ... Get the extension and content type of file
	GetExtensionAndContentType(file io.Reader) (string, string, error)

	// ValidateFileType ... Validate file by its required type
	ValidateFileType(extension string) error
}

func Bootstrap(provider string) {
	dm := dependency.GetManager()
	var fm interface{}
	switch provider {
	case ProviderAwsS3:
		fm = &awsS3Manager{
			cm:  dm.Get(config.DependencyConfigManager).(config.Manager),
			clm: dm.Get(cloud.DependencyCloudManager).(cloud.Manager),
		}
		fm.(*awsS3Manager).initialize()
	}
	dm.Register(DependencyFileManager, fm)
}

func getImageSupportedExtensions() map[string]string {
	return map[string]string{"jpg": "jpg", "jpeg": "jpeg", "png": "png", "jfif": "jfif", "pjpeg": "pjpeg", "pjp": "pjp"}
}
