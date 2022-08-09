package exception

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ErrorManagerTestSuite struct {
	suite.Suite
}

func TestErrorManager(t *testing.T) {
	suite.Run(t, new(ErrorManagerTestSuite))
}

func (ts *ErrorManagerTestSuite) Test_exceptionManager_UserException() {
	ts.Run(
		"success", func() {
			err := CreateUserException("user already exists")
			assert.NotNil(ts.T(), err.Error())
		},
	)
	ts.Run(
		"success - with data", func() {
			err := CreateUserException("user does not exist", "jon")
			assert.NotNil(ts.T(), err.Error())
		},
	)
}

func (ts *ErrorManagerTestSuite) Test_exceptionManager_AppException() {
	ts.Run(
		"success", func() {
			err := CreateAppException("internal server exception")
			assert.NotNil(ts.T(), err.Error())
		},
	)
	ts.Run(
		"success - with data", func() {
			err := CreateAppException("internal server exception", "jon")
			assert.NotNil(ts.T(), err.Error())
		},
	)
}

func (ts *ErrorManagerTestSuite) Test_exceptionManager_UnauthorizedException() {
	ts.Run(
		"success", func() {
			err := CreateUnauthorizedException("unauthorized")
			assert.NotNil(ts.T(), err.Error())
		},
	)
	ts.Run(
		"success - with data", func() {
			err := CreateUnauthorizedException("unauthorized", "jon")
			assert.NotNil(ts.T(), err.Error())
		},
	)
}

func (ts *ErrorManagerTestSuite) Test_exceptionManager_AccessDeniedException() {
	ts.Run(
		"success", func() {
			err := CreateAccessDeniedException("access denied")
			assert.NotNil(ts.T(), err.Error())
		},
	)
	ts.Run(
		"success - with data", func() {
			err := CreateAccessDeniedException("access denied", "jon")
			assert.NotNil(ts.T(), err.Error())
		},
	)
}

func (ts *ErrorManagerTestSuite) Test_exceptionManager_ValidationException() {
	ts.Run(
		"success", func() {
			err := CreateValidationException("validation failed")
			assert.NotNil(ts.T(), err.Error())
		},
	)

	ts.Run(
		"success - with data", func() {
			err := CreateValidationException("validation failed", "invalid mobile number")
			assert.NotNil(ts.T(), err.Error())
		},
	)
}
