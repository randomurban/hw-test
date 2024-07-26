package main

import (
	"errors"
	"io"
	"os"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	fromFile, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer fromFile.Close()

	fromInfo, err := fromFile.Stat()
	if err != nil {
		return err
	}

	if offset < 0 || offset > fromInfo.Size() {
		return ErrOffsetExceedsFileSize
	}
	_, err = fromFile.Seek(offset, io.SeekStart)
	if err != nil {
		return err
	}

	toFile, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer toFile.Close()

	limitTo := limit
	if limitTo == 0 || limitTo > fromInfo.Size()-offset {
		limitTo = fromInfo.Size() - offset
	}
	_, err = io.CopyN(toFile, fromFile, limitTo)
	if err != nil {
		return err
	}

	return toFile.Sync()
}
