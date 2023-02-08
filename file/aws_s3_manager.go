package file

import (
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gabriel-vasile/mimetype"

	"github.com/adwitiyaio/arka/cloud"
	"github.com/adwitiyaio/arka/config"
	"github.com/adwitiyaio/arka/exception"
	"github.com/adwitiyaio/arka/logger"
)

const regionKey = "AWS_REGION"
const storageBucketKey = "AWS_STORAGE_BUCKET"

type awsS3Manager struct {
	cm     config.Manager
	clm    cloud.Manager
	client *s3.Client
	region string
	bucket string
}

func (cfm *awsS3Manager) initialize() {
	config := cfm.clm.GetConfig()
	cfm.client = s3.NewFromConfig(config)
	cfm.bucket = cfm.cm.GetValueForKey(storageBucketKey)
	cfm.region = cfm.cm.GetValueForKey(regionKey)
}

func (cfm *awsS3Manager) UploadFile(filename string, contentType string, file io.Reader, directory *string) (string, error) {
	var key, url string
	if directory != nil {
		key = fmt.Sprintf("%s/%s", *directory, filename)
	} else {
		key = filename
	}

	upParams := &s3.PutObjectInput{
		Bucket:      aws.String(cfm.bucket),
		Key:         aws.String(key),
		Body:        file,
		ContentType: aws.String(contentType),
	}
	_, err := cfm.client.PutObject(context.TODO(), upParams)
	if err != nil {
		logger.Log.Error().Err(err).Msg("failed to upload file to S3")
		return url, err
	}
	url = cfm.constructUrlForFile(key)
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
