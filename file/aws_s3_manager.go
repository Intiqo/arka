package file

import (
	"fmt"
	"github.com/adwitiyaio/arka/cloud"
	"github.com/adwitiyaio/arka/config"
	"github.com/adwitiyaio/arka/exception"
	"github.com/adwitiyaio/arka/logger"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gabriel-vasile/mimetype"
	"io"
)

const regionKey = "AWS_REGION"
const storageBucketKey = "AWS_STORAGE_BUCKET"

type awsS3Manager struct {
	cm      config.Manager
	clm     cloud.Manager
	session *session.Session
	region  string
	bucket  string
}

func (cfm *awsS3Manager) initialize() {
	cfm.bucket = cfm.cm.GetValueForKey(storageBucketKey)
	cfm.session = cfm.clm.GetSession()
	cfm.region = cfm.cm.GetValueForKey(regionKey)
}

func (cfm *awsS3Manager) UploadFile(filename string, contentType string, file io.Reader, directory *string) (string, error) {
	uploader := s3manager.NewUploader(cfm.session)
	var key, url string
	if directory != nil {
		key = fmt.Sprintf("%s/%s", *directory, filename)
	} else {
		key = filename
	}
	upParams := &s3manager.UploadInput{
		Bucket:      aws.String(cfm.bucket),
		Key:         aws.String(key),
		Body:        file,
		ContentType: aws.String(contentType),
	}
	result, err := uploader.Upload(upParams)
	if err != nil {
		logger.Log.Error().Err(err).Msg("failed to upload file to S3")
		return url, err
	}
	url = cfm.constructUrlForFile(key)

	logger.Log.Debug().Str("name", filename).Str("location", result.Location).Msgf("uploaded file to s3 location")
	return url, nil
}

func (cfm *awsS3Manager) GetExtensionAndContentType(file io.Reader) (string, string, error) {
	var extension string
	var contentType string
	mimeType, err := mimetype.DetectReader(file)

	if err != nil {
		return "", "", err
	} else {
		extension = mimeType.Extension()
		contentType = mimeType.String()
	}
	return extension, contentType, nil
}

func (cfm *awsS3Manager) ValidateFileType(extension string) error {
	extensionMap := getImageSupportedExtensions()
	_, ok := extensionMap[extension[1:]]
	if !ok {
		return exception.CreateAppException(ErrorUnsupportedFileType)
	}
	return nil
}

func (cfm *awsS3Manager) constructUrlForFile(key string) string {
	return fmt.Sprintf("https://%s.s3-%s.amazonaws.com/%s", cfm.bucket, cfm.region, key)
}
