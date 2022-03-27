package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func setUp() (dirname string) {
	var err error
	dirname, err = os.MkdirTemp("", "*")
	if err != nil {
		panic(err)
	}

	fmt.Println("temp dirname created: ", dirname)

	return
}

func tearDown(dirname string) {
	fmt.Println("temp dirname removed: ", dirname)
	os.RemoveAll(dirname)
}

func TestCopy(t *testing.T) {
	t.Run("unsupported file", func(t *testing.T) {
		testDir := setUp()
		defer tearDown(testDir)

		dstFileName := filepath.Join(testDir, "file")

		err := Copy(".", dstFileName, 0, 0)
		require.Equal(t, err, ErrUnsupportedFile)
		err = Copy("/dev/urandom", dstFileName, 0, 0)
		require.Equal(t, err, ErrUnsupportedFile)
	})

	t.Run("file not found", func(t *testing.T) {
		testDir := setUp()
		defer tearDown(testDir)

		err := Copy(filepath.Join(testDir, "NOT_EXISTED_FILENAME"), filepath.Join(testDir, "file"), 0, 0)
		require.Equal(t, err, ErrSrcFileNotFound)
	})

	t.Run("file is empty", func(t *testing.T) {
		testDir := setUp()
		defer tearDown(testDir)

		emptyFileName := filepath.Join(testDir, "NOT_EXISTED_FILENAME")

		_, err := os.Create(emptyFileName)
		if err != nil {
			panic(err)
		}
		err = Copy(emptyFileName, filepath.Join(testDir, "file"), 0, 0)
		require.Equal(t, err, ErrFileIsEmpty)
	})

	t.Run("offset violation", func(t *testing.T) {
		testDir := setUp()
		defer tearDown(testDir)

		fi, err := os.Stat("testdata/input.txt")
		if err != nil {
			panic(err)
		}

		outFileName := filepath.Join(testDir, "out")

		err = Copy("testdata/input.txt", outFileName, fi.Size(), 0)
		require.Equal(t, err, ErrOffsetExceedsFileSize)

		err = Copy("testdata/input.txt", outFileName, fi.Size()+1, 0)
		require.Equal(t, err, ErrOffsetExceedsFileSize)
	})
}
