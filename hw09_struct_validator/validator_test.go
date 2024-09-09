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
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
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
	WrongTag struct {
		Code int    `validate:"len:100,500"`
		Body string `validate:"omitempty"`
	}
	WrongRegexp struct {
		Code int    `validate:"min:100"`
		Body string `validate:"regexp:a(b"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
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
			in:          AppWrong{"ver_1"},
			expectedErr: errors.New("parsing AppWrong: invalid number param for len: strconv.Atoi: parsing \"five\": invalid syntax"),
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
			in:          AppIntWrong{[]int{1, 2}},
			expectedErr: errors.New("parsing AppIntWrong: invalid int tag element: ver.1"),
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
		{
			in: WrongTag{
				100500,
				"body",
			},
			expectedErr: errors.New("parsing WrongTag: invalid number param for len: strconv.Atoi: parsing \"100,500\": invalid syntax"),
		},
		{
			in: WrongRegexp{
				Code: 0,
				Body: "bad regexp",
			},
			expectedErr: errors.New("parsing WrongRegexp: error parsing regexp: missing closing ): `a(b`"),
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
				if !errors.As(got, &validationErr) && got.Error() != expErr {
					t.Errorf("got %q, wanted %v", got, tt.expectedErr)
				}
			}
		})
	}
}

func TestValidate2(t *testing.T) {
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
				if !errors.As(got, &validationErr) && got.Error() != expErr {
					t.Errorf("got %q, wanted %v", got, tt.expectedErr)
				}
			}
		})
	}
}
