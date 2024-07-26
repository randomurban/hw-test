package main

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestCase struct {
	input  string
	output string
	offset int64
	limit  int64
}

func TestCopy(t *testing.T) {
	testCases := []TestCase{
		{"testdata/input.txt", "testdata/out_offset6000_limit1000.txt", 6000, 1_000},
		{"testdata/input.txt", "testdata/out_offset0_limit0.txt", 0, 0},
		{"testdata/input.txt", "testdata/out_offset0_limit10.txt", 0, 10},
		{"testdata/input.txt", "testdata/out_offset0_limit1000.txt", 0, 1_000},
		{"testdata/input.txt", "testdata/out_offset0_limit10000.txt", 0, 10_000},
		{"testdata/input.txt", "testdata/out_offset100_limit1000.txt", 100, 1_000},
	}

	for _, tc := range testCases {
		input := tc.input
		output := tc.output
		offset := tc.offset
		limit := tc.limit
		outDir, outName := path.Split(output)
		result := outDir + "res_" + outName

		err := Copy(input, result, offset, limit)
		if err != nil {
			t.Error(err.Error())
		}
	}
	t.Run("check size of copy", func(t *testing.T) {
		for _, tc := range testCases {
			reference := tc.output
			outDir, outName := path.Split(reference)
			result := outDir + "res_" + outName

			refFile, err := os.Open(reference)
			if err != nil {
				t.Error(err.Error())
			}
			defer refFile.Close()

			refStat, err := refFile.Stat()
			if err != nil {
				t.Error(err.Error())
			}

			resultFile, err := os.Open(result)
			if err != nil {
				t.Error(err.Error())
			}
			defer resultFile.Close()
			resultStat, err := resultFile.Stat()
			if err != nil {
				t.Error(err.Error())
			}

			assert.Equal(t, refStat.Size(), resultStat.Size())
		}
	})
	for _, tc := range testCases {
		outDir, outName := path.Split(tc.output)
		result := outDir + "res_" + outName
		err := os.Remove(result)
		if err != nil {
			t.Error(err.Error())
		}
	}
}
