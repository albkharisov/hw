package main

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
)

type EnvTestSuite struct {
	suite.Suite
	testDir string
}

func (s *EnvTestSuite) SetupTest() {
	var err error
	s.testDir, err = os.MkdirTemp("", "*")
	if err != nil {
		panic(err)
	}
}

func (s *EnvTestSuite) TearDownTest() {
	os.RemoveAll(s.testDir)
}

func TestReadDir(t *testing.T) {
	suite.Run(t, new(EnvTestSuite))
}

func (s *EnvTestSuite) TestDirNotExists() {
	// create temp directory to be sure it's empty
	_, err := ReadDir(filepath.Join(s.testDir, "not_existed_directory"))

	s.Require().NotNil(err) //nolint:typecheck
	var ve *os.PathError
	s.Require().ErrorAs(err, &ve) //nolint:typecheck
}

func createTempDir(path string) (outName string) {
	outName, err := os.MkdirTemp(path, "*")
	if err != nil {
		panic(err)
	}

	return
}

func fillString(size int64) string {
	str := make([]rune, size)
	for i := int64(0); i < size; i++ {
		str[i] = 'a'
	}

	return string(str)
}

func createTempFile(path string, size int64) {
	f, err := os.CreateTemp(path, "*")
	if err != nil {
		panic(err)
	}

	reader := strings.NewReader(fillString(size))

	_, err = io.CopyN(f, reader, size)
	if err != nil {
		panic(err)
	}
}

func (s *EnvTestSuite) TestAllFilesFound() {
	var count int

	dir1 := createTempDir(s.testDir)
	dir2 := createTempDir(s.testDir)

	createTempFile(s.testDir, 0)
	count++
	createTempFile(dir1, 1)
	count++
	createTempFile(dir1, 513)
	createTempFile(dir1, 512)
	count++
	createTempFile(dir2, 100)
	count++

	envmap, err := ReadDir(s.testDir)
	s.Require().NoError(err)              //nolint:typecheck
	s.Require().Equal(count, len(envmap)) //nolint:typecheck
}
