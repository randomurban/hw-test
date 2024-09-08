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

	AppIn struct {
		Version []string `validate:"in:ver.1,ver_1"`
	}

	AppIntWrong struct {
		Version []int `validate:"in:ver.1,1"`
	}

	AppInt struct {
		Version []int `validate:"in:1,2"`
	}
	AppMin struct {
		Version []int `validate:"min:10"`
	}

	AppMax struct {
		Version []int `validate:"max:20"`
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
			in: User{
				ID:     "123456789012345678901234567890123456",
				Name:   "*skipped*",
				Age:    18,
				Email:  "e@mail.com",
				Role:   "admin",
				Phones: []string{"12345678901", ""},
				meta:   json.RawMessage("{}"),
			},
			expectedErr: ValidationErrors{
				{
					"Phones",
					errors.New("len must be 11"),
				},
			},
		},

		{
			in: User{
				ID:     "123456789012345678901234567890123456",
				Name:   "*skipped*",
				Age:    18,
				Email:  "email",
				Role:   "admin",
				Phones: nil,
				meta:   nil,
			},
			expectedErr: ValidationErrors{{
				"Email",
				errors.New("'email' not matched"),
			}},
		},
		{
			in: AppSlice{[]string{""}},
			expectedErr: ValidationErrors{{
				"Version",
				errors.New("len must be 5"),
			}},
		},
		{
			in:          AppSlice{[]string{}},
			expectedErr: nil,
		},
		{
			in:          Token{[]byte{}, []byte{}, []byte{}},
			expectedErr: nil,
		},
		{
			in:          Response{200, "empty"},
			expectedErr: nil,
		},
		{
			in:          Response{100, "empty"},
			expectedErr: ValidationErrors{{"Code", errors.New("100 in [200,404,500] is required")}},
		},
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
		{
			in: AppIn{[]string{"ver1.0", "ver.1"}},
			expectedErr: ValidationErrors{{
				"Version",
				errors.New("ver1.0 in [ver.1,ver_1] is required"),
			}},
		},
		{
			in: AppIntWrong{[]int{1, 2}},
			expectedErr: ValidationErrors{
				{
					"Version",
					errors.New("invalid int tag element: ver.1"),
				},
				{
					"Version",
					errors.New("invalid int tag element: ver.1"),
				},
			},
		},
		{
			in:          AppInt{[]int{1, 2}},
			expectedErr: nil,
		},
		{
			in: AppInt{[]int{1, 2, 3}},
			expectedErr: ValidationErrors{
				{
					"Version",
					errors.New("3 in [1,2] is required"),
				},
			},
		},

		{
			in: AppMin{[]int{1, 10}},
			expectedErr: ValidationErrors{{
				"Version",
				errors.New("must be minimum 10"),
			}},
		},
		{
			in: AppMax{[]int{10, 30}},
			expectedErr: ValidationErrors{{
				"Version",
				errors.New("must be maximum 20"),
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
						t.Errorf("got %q, wanted %v", got, tt.expectedErr)
					}
				} else {
					t.Errorf("got %v, wanted %v", got, tt.expectedErr)
				}
			} else {
				if tt.expectedErr != nil {
					t.Errorf("got %v, wanted %v", got, tt.expectedErr)
				}
			}
		})
	}
}
