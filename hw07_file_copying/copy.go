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

	fromInfo, err := fromFile.Stat()
	if err != nil {
		return err
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

	toFile, err := os.Create(toPath)
	if err != nil {
		return err
	}

	barWriter := bar.NewProxyWriter(toFile)

	_, err = io.Copy(barWriter, limitReader)
	if err != nil {
		return err
	}

	err = fromFile.Close()
	if err != nil {
		return err
	}
	err = toFile.Sync()
	if err != nil {
		return err
	}
	bar.Finish()
	println()

	err = toFile.Close()
	return err
}
