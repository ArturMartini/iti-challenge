package service

import (
	"github.com/arturmartini/iti-challenge/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type passwordSuiteTest struct {
	suite.Suite
}

type testCase struct {
	Value    string
	Expected bool
}

func (*passwordSuiteTest) SetupTest() {
	instance = service{}
}

func TestRunPassword(t *testing.T) {
	suite.Run(t, new(passwordSuiteTest))
}

func (suite *passwordSuiteTest) TestPasswordValidate() {
	testCases := []testCase{
		{
			Value:    "",
			Expected: false,
		},
		{
			Value:    "aa",
			Expected: false,
		},
		{
			Value:    "ab",
			Expected: false,
		},
		{
			Value:    "AAAbbbCc",
			Expected: false,
		},
		{
			Value:    "AbTp9!foo",
			Expected: false,
		},
		{
			Value:    "AbTp9!foA",
			Expected: false,
		},
		{
			Value:    "AbTp9 fok",
			Expected: false,
		},
		{
			Value:    "AbTp9!fok",
			Expected: true,
		},
	}

	validateTest(suite, testCases)
}

func validateTest(suite *passwordSuiteTest, testCases []testCase) {
	for _, c := range testCases {
		pass := entities.Password{Value: c.Value}
		assert.Equal(suite.T(), c.Expected, instance.ValidateStrongPassword(pass))
	}
}
