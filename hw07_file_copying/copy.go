package main

import (
	"errors"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, limit, offset int64) error {
	fromFile, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer fromFile.Close()

	fromInfo, err := fromFile.Stat()
	if err != nil {
		return err
	}
	if !fromInfo.Mode().IsRegular() {
		return ErrUnsupportedFile
	}

	fromSize := fromInfo.Size()

	if offset < 0 || offset > fromSize {
		return ErrOffsetExceedsFileSize
	}
	_, err = fromFile.Seek(offset, io.SeekStart)
	if err != nil {
		return err
	}

	limitTo := limit
	if limit == 0 || limit > fromSize-offset {
		limitTo = fromSize - offset
	}

	limitReader := io.LimitReader(fromFile, limitTo)
	bar := pb.Simple.Start64(limitTo)

	toFile, err := os.CreateTemp("", "*.tmp")
	if err != nil {
		return err
	}
	defer toFile.Close()
	println(toFile.Name())

	barWriter := bar.NewProxyWriter(toFile)

	_, err = io.Copy(barWriter, limitReader)
	if err != nil {
		return err
	}

	err = toFile.Close()
	if err != nil {
		return err
	}

	bar.Finish()
	println()

	err = os.Rename(toFile.Name(), toPath)
	return err
}
