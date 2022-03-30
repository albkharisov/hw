package main

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type RunCmdTestSuite struct {
	suite.Suite
}

func (s *RunCmdTestSuite) SetupTest() {
}

func (s *RunCmdTestSuite) TearDownTest() {
}

func TestRunCmd(t *testing.T) {
	suite.Run(t, new(RunCmdTestSuite))
}

func (s *RunCmdTestSuite) TestExitResult() {
	result := RunCmd([]string{"false"}, Environment{})
	s.Require().Equal(1, result)

	result = RunCmd([]string{"true"}, Environment{})
	s.Require().Equal(0, result)
}
