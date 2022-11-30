package util

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/adwitiyaio/arka/dependency"
)

type AppUtilTestSuite struct {
	suite.Suite

	a AppUtil
}

func (ts *AppUtilTestSuite) SetupSuite() {
	Bootstrap()
	ts.a = dependency.GetManager().Get(DependencyAppUtil).(AppUtil)
}

func TestAppUtil(t *testing.T) {
	suite.Run(t, new(AppUtilTestSuite))
}

func (ts *AppUtilTestSuite) Test_simpleAppUtil_GetCurrentTime() {
	ts.Run(
		"success", func() {
			t := ts.a.GetCurrentTime()
			assert.NotNil(ts.T(), t)
		},
	)
}

func (ts *AppUtilTestSuite) Test_OtpGenerator() {
	ts.Run(
		"should generate 2 digits otp", func() {
			t := ts.a.GenerateOTP(2)
			assert.NotNil(ts.T(), t)
			assert.Equal(ts.T(), 2, len(t))
		},
	)

	ts.Run(
		"should generate 6 digit", func() {
			t := ts.a.GenerateOTP(6)
			assert.NotNil(ts.T(), t)
			assert.Equal(ts.T(), 6, len(t))
		},
	)

	ts.Run(
		"should not generate otp when the length is 0", func() {
			t := ts.a.GenerateOTP(0)
			assert.NotNil(ts.T(), t)
			assert.Equal(ts.T(), "", t)
		},
	)

	ts.Run(
		"should not generate otp when the length is -1", func() {
			t := ts.a.GenerateOTP(-1)
			assert.NotNil(ts.T(), t)
			assert.Equal(ts.T(), "", t)
		},
	)
}

func (ts *AppUtilTestSuite) Test_simpleAppUtil_GenerateUniqueToken() {
	ts.Run(
		"success", func() {
			t := ts.a.GenerateUniqueToken()
			assert.NotNil(ts.T(), t)
		},
	)
}

func (ts *AppUtilTestSuite) Test_simpleAppUtil_GetExpiryTimeForDuration() {
	ts.Run(
		"success - future", func() {
			t := ts.a.GetExpiryTimeForDuration(24)
			assert.True(ts.T(), time.Now().Before(t))
		},
	)

	ts.Run(
		"success - past", func() {
			t := ts.a.GetExpiryTimeForDuration(-24)
			assert.True(ts.T(), time.Now().After(t))
		},
	)
}

func (ts *AppUtilTestSuite) Test_simpleAppUtil_ParseStringForTime() {
	ts.Run(
		"invalid iso time string", func() {
			_, err := ts.a.ParseStringForTime("2021-06-16T14:90:00Z")
			require.Error(ts.T(), err)
		},
	)

	ts.Run(
		"success", func() {
			t, err := ts.a.ParseStringForTime("2021-06-16T14:30:00Z")
			require.NoError(ts.T(), err)
			assert.NotNil(ts.T(), t)
		},
	)
}

func (ts *AppUtilTestSuite) Test_simpleAppUtil_ParseStringForTimeWithLocation() {
	ts.Run(
		"invalid iso time string", func() {
			_, err := ts.a.ParseStringForTime("2021-06-16T14:90:00Z")
			require.Error(ts.T(), err)
		},
	)

	ts.Run(
		"success", func() {
			loc, err := time.LoadLocation("Asia/Kolkata")
			require.NoError(ts.T(), err)
			t, err := ts.a.ParseStringForTimeWithLocation("2021-06-16T14:30:00+05:30", loc)
			require.NoError(ts.T(), err)
			assert.NotNil(ts.T(), t)
		},
	)
}

func (ts *AppUtilTestSuite) Test_simpleAppUtil_CompareSlices() {
	ts.Run(
		"success", func() {
			a := []string{"abc", "def"}
			b := []string{"abc"}
			result := ts.a.CompareSlices(a, b)
			assert.NotNil(ts.T(), result)
			assert.Equal(ts.T(), "def", result[0])
		},
	)
}

func (ts *AppUtilTestSuite) Test_simpleAppUtil_IsTimeExpired() {
	ts.Run(
		"expired", func() {
			t := time.Now()
			newT := t.Add(-time.Hour * 1)

			result := ts.a.IsTimeExpired(newT)
			assert.NotNil(ts.T(), result)
			assert.Equal(ts.T(), true, result)
		},
	)

	ts.Run(
		"not expired", func() {
			t := time.Now()
			newT := t.Add(time.Hour * 1)

			result := ts.a.IsTimeExpired(newT)
			assert.NotNil(ts.T(), result)
			assert.Equal(ts.T(), false, result)
		},
	)
}

func (ts *AppUtilTestSuite) Test_simpleAppUtil_FormatDate() {
	ts.Run(
		"success", func() {
			date, err := time.Parse(time.RFC3339, "2021-06-15T00:00:00Z")
			assert.NoError(ts.T(), err)
			result := ts.a.FormatDate(date)
			assert.Equal(ts.T(), "15th Jun", result)
		},
	)
}

func (ts *AppUtilTestSuite) Test_simpleAppUtil_ParseWeekday() {
	ts.Run(
		"invalid weekday", func() {
			_, err := ts.a.ParseWeekday("November")
			assert.Error(ts.T(), err)
			assert.Equal(ts.T(), "invalid weekday 'November'", err.Error())
		},
	)

	ts.Run(
		"success", func() {
			result, err := ts.a.ParseWeekday(time.Monday.String())
			assert.NoError(ts.T(), err)
			assert.Equal(ts.T(), time.Monday, result)
		},
	)
}
