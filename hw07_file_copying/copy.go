package main

import (
	"io"
	"os"

	"github.com/cheggaaa/pb/v3" //nolint:typecheck
	"github.com/pkg/errors"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrSrcFileNotFound       = errors.New("source file not found")
	ErrFileIsEmpty           = errors.New("source file is empty")
)

func min(a, b int64) int64 {
	if a > b {
		return b
	}
	return a
}

func Copy(fromPath, toPath string, offset, limit int64) error {
	// check we don't write into irregular file
	fi, err := os.Stat(toPath)
	if err == nil {
		if !fi.Mode().IsRegular() {
			return ErrUnsupportedFile
		}
	}

	// check source file exists
	fi, err = os.Stat(fromPath)
	if err != nil {
		return ErrSrcFileNotFound
	}

	// check source file can be opened
	fileFrom, err := os.Open(fromPath)
	if err != nil {
		return errors.Wrap(err, "Open source failed!")
	}
	defer fileFrom.Close()

	// check destination file can be opened/created
	fileTo, err := os.OpenFile(toPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return errors.Wrap(err, "OpenFile destination failed!")
	}
	defer fileTo.Close()

	// check source file is regular
	if !fi.Mode().IsRegular() {
		return ErrUnsupportedFile
	}

	// check source file is not empty
	if fi.Size() == 0 {
		return ErrFileIsEmpty
	}

	// check offset is less than filesize
	if offset >= fi.Size() {
		return ErrOffsetExceedsFileSize
	}

	var copySize int64
	if limit == 0 {
		copySize = fi.Size() - offset
	} else {
		copySize = min(fi.Size()-offset, limit)
	}
	progbar := pb.New64(copySize) //nolint:typecheck
	fileFrom.Seek(offset, io.SeekStart)
	proxyReader := progbar.NewProxyReader(fileFrom)

	progbar.Start()
	_, err = io.CopyN(fileTo, proxyReader, copySize)
	progbar.Finish()

	if errors.Is(err, io.EOF) {
		// https://www.joeshaw.org/dont-defer-close-on-writable-files/
		return fileTo.Sync()
	}
	return err
}
