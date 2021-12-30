package file

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/adwitiyaio/arka/cloud"
	"github.com/adwitiyaio/arka/config"
	"github.com/adwitiyaio/arka/dependency"
)

type FileManagerTestSuite struct {
	suite.Suite
	flm Manager
}

func TestFileManager(t *testing.T) {
	suite.Run(t, new(FileManagerTestSuite))
}

func (ts *FileManagerTestSuite) SetupSuite() {
	dm := dependency.GetManager()
	config.Bootstrap(config.ProviderEnvironment, "../test.env")
	cloud.Bootstrap(cloud.ProviderAws)
	Bootstrap(ProviderAwsS3)
	ts.flm = dm.Get(DependencyFileManager).(Manager)
}

func (ts *FileManagerTestSuite) Test_awsFileManager_UploadFile() {
	ts.Run("unsupported file type", func() {
		data, err := ioutil.ReadFile("./testdata/sample.txt")
		require.NoError(ts.T(), err)

		file := bytes.NewReader(data)
		require.NotNil(ts.T(), file)

		extension, _, err := ts.flm.GetExtensionAndContentType(file)
		require.NoError(ts.T(), err)

		err = ts.flm.ValidateFileType(extension)
		require.Error(ts.T(), err)
	})

	ts.Run("success - upload file", func() {
		data, err := ioutil.ReadFile("./testdata/sample.png")
		require.NoError(ts.T(), err)

		file := bytes.NewReader(data)
		require.NotNil(ts.T(), file)

		extension, contentType, err := ts.flm.GetExtensionAndContentType(file)
		require.NoError(ts.T(), err)

		err = ts.flm.ValidateFileType(extension)
		require.NoError(ts.T(), err)

		filename := "sample" + extension
		url, err := ts.flm.UploadFile(filename, contentType, file, nil)
		assert.NoError(ts.T(), err)
		assert.NotNil(ts.T(), url)
	})

	ts.Run("success - upload file with directory", func() {
		data, err := ioutil.ReadFile("./testdata/sample.png")
		require.NoError(ts.T(), err)

		file := bytes.NewReader(data)
		require.NotNil(ts.T(), file)

		extension, contentType, err := ts.flm.GetExtensionAndContentType(file)
		require.NoError(ts.T(), err)

		err = ts.flm.ValidateFileType(extension)
		require.NoError(ts.T(), err)

		filename := "sample" + extension
		directory := "logo"
		url, err := ts.flm.UploadFile(filename, contentType, file, &directory)
		assert.NoError(ts.T(), err)
		assert.NotNil(ts.T(), url)
	})
}
