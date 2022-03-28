package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"
)

type CopyTestSuite struct {
	suite.Suite
	testDir string
}

func (s *CopyTestSuite) SetupTest() {
	var err error
	s.testDir, err = os.MkdirTemp("", "*")
	if err != nil {
		panic(err)
	}
}

func (s *CopyTestSuite) TearDownTest() {
	os.RemoveAll(s.testDir)
}

func (s *CopyTestSuite) TestUnsupportedFile() {
	srcFileName := "/dev/urandom"
	dstFileName := filepath.Join(s.testDir, "file")

	err := Copy(".", dstFileName, 0, 0)
	s.Require().Equal(err, ErrUnsupportedFile)

	err = Copy(srcFileName, dstFileName, 0, 0)
	s.Require().Equal(err, ErrUnsupportedFile)
}

func (s *CopyTestSuite) TestFileNotFound() {
	srcFileName := filepath.Join(s.testDir, "NOT_EXISTED_FILENAME")
	dstFileName := filepath.Join(s.testDir, "file")

	err := Copy(srcFileName, dstFileName, 0, 0)
	s.Require().Equal(err, ErrSrcFileNotFound)
}

func (s *CopyTestSuite) TestFileIsEmpty() {
	srcFileName := filepath.Join(s.testDir, "NOT_EXISTED_FILENAME")
	dstFileName := filepath.Join(s.testDir, "file")

	_, err := os.Create(srcFileName)
	if err != nil {
		panic(err)
	}

	err = Copy(srcFileName, dstFileName, 0, 0)
	s.Require().Equal(err, ErrFileIsEmpty)
}

func (s *CopyTestSuite) TestOffsetViolation() {
	srcFileName := "testdata/input.txt"
	dstFileName := filepath.Join(s.testDir, "out")

	fi, err := os.Stat(srcFileName)
	if err != nil {
		panic(err)
	}

	err = Copy(srcFileName, dstFileName, fi.Size(), 0)
	s.Require().Equal(err, ErrOffsetExceedsFileSize)

	err = Copy(srcFileName, dstFileName, fi.Size()+1, 0)
	s.Require().Equal(err, ErrOffsetExceedsFileSize)
}

func TestCopy(t *testing.T) {
	suite.Run(t, new(CopyTestSuite))
}
