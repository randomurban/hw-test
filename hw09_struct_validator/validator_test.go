package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
	}

	App struct {
		Version string `validate:"len:5"`
	}

	AppWrong struct {
		Version string `validate:"len:five"`
	}

	AppSlice struct {
		Version []string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in:          App{"ver1.0"},
			expectedErr: ValidationErrors{{"Version", errors.New("len must be 5")}},
		},
		{
			in:          App{"ver_1"},
			expectedErr: nil,
		},
		{
			in: AppWrong{"ver_1"},
			expectedErr: ValidationErrors{
				{
					"Version",
					errors.New("parsing AppWrong: number param for len: strconv.Atoi: parsing \"five\": invalid syntax"),
				},
			},
		},
		{
			in: AppSlice{[]string{"ver1.0", "ver2.0"}},
			expectedErr: ValidationErrors{{
				"Version",
				errors.New("len must be 5"),
			}, {"Version", errors.New("len must be 5")}},
		},
		{
			in: AppSlice{[]string{"ver1.0", "ver.2"}},
			expectedErr: ValidationErrors{{
				"Version",
				errors.New("len must be 5"),
			}},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()
			var expErr string
			if tt.expectedErr != nil {
				expErr = tt.expectedErr.Error()
			}
			got := Validate(tt.in)
			var validationErr ValidationErrors
			if got != nil {
				if errors.As(got, &validationErr) {
					if got.Error() != expErr {
						t.Errorf("got %v, wanted %v", got, tt.expectedErr)
					}
				} else {
					t.Errorf("got %v, wanted %v", got, tt.expectedErr)
				}
			}
		})
	}
}
