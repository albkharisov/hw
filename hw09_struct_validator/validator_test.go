package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
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
			in:          App{Version: "12345"},
			expectedErr: nil,
		},
		{
			in:          App{Version: "123456"},
			expectedErr: ValidationErrors{{Field: "Version", Err: LenError}},
		},
		{
			in:          App{Version: ""},
			expectedErr: ValidationErrors{{Field: "Version", Err: LenError}},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)

			if tt.expectedErr == nil {
				fmt.Println("expected == nil, return, Real err: ", err, ", version: ", tt.in.(App).Version)
				require.Nil(t, err)
				return
			}

			switch expErrs := tt.expectedErr.(type) {
			case ValidationErrors:
				resultErrs, ok := err.(ValidationErrors)
				require.True(t, ok, "Expected ValidationErrors")
				require.Equal(t, len(expErrs), len(resultErrs))

				for j := range expErrs {
					require.Equal(t, expErrs[j].Err, resultErrs[j].Err)
				}
			default:
				require.Equal(t, tt.expectedErr, err)
			}
		})
	}
}
