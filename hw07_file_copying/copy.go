package main

import (
	"errors"
	"io"
	"os"
	"time"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

const BufSize = 64

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

	reader := io.LimitReader(fromFile, limitTo)
	bar := pb.Simple.Start64(limitTo).SetRefreshRate(time.Millisecond * 10)
	bar.Start()
	barReader := bar.NewProxyReader(reader)

	toFile, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer toFile.Close()

	_, err = io.Copy(toFile, barReader)
	if err != nil {
		return err
	}

	bar.Finish()
	println()
	return toFile.Sync()
}
