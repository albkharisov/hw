package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:6"`
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

	Malformed1 struct {
		Field1 string `validate:"hz:12"`
	}

	Malformed2 struct {
		Field1 string `validate:""`
	}

	Malformed3 struct {
		Field1 string `validate:"regexp:^\\w+@\\w++\\.\\w+$"`
	}

	Malformed4 struct {
		Field1 string `validate:"min:12"`
	}

	MultiFailed1 struct {
		Field1 int `validate:"min:12|min:10|min:11|max:1|in:9,2,3,4,1,9"`
	}

	StructTest1 struct {
		Field1 int         `validate:"min:12"`
		Field2 StructTest2 `validate:"nested"`
	}

	StructTest2 struct {
		Field2_1 StructTest3 `validate:"nested"`
		Field2_2 int         `validate:"max:10"`
	}

	StructTest3 struct {
		Field3_1 int          `validate:"in:1,2"`
		Field3_2 string       `validate:"in:3,4"`
		Field3_3 MultiFailed1 // no 'nested' tag - no validation
		Field3_4 Malformed1   // no 'nested' tag - no validation
	}

	StructTest4 struct {
		Field4_1 int        `validate:"in:2"`
		Field4_2 Malformed1 `validate:"nested"`
	}
)

var tests = []struct {
	in          interface{}
	expectedErr error
}{
	{
		in: User{
			ID:     "123456",          // string `json:"id" validate:"len:6"`
			Age:    20,                // int      `validate:"min:18|max:50"`
			Email:  "baton@hleba.com", // string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
			Role:   "admin",           // UserRole `validate:"in:admin,stuff"`
			Phones: nil,               // []string `validate:"len:11"`
			meta:   json.RawMessage{},
		},
		expectedErr: nil,
	},
	{
		in: User{
			ID:     "12345",                                // string `json:"id" validate:"len:6"`
			Age:    20,                                     // int      `validate:"min:18|max:50"`
			Email:  "baton@hleba.com",                      // string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
			Role:   "admin",                                // UserRole `validate:"in:admin,stuff"`
			Phones: []string{"89125885738", "89125885736"}, // []string `validate:"len:11"`
		},
		expectedErr: ValidationErrors{{Field: "ID", Err: ErrLen}},
	},
	{
		in: User{
			ID:     "123456",                               // string `json:"id" validate:"len:6"`
			Age:    2,                                      // int      `validate:"min:18|max:50"`
			Email:  "baton@hleba.com",                      // string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
			Role:   "admin",                                // UserRole `validate:"in:admin,stuff"`
			Phones: []string{"89125885738", "89125885736"}, // []string `validate:"len:11"`
		},
		expectedErr: ValidationErrors{{Field: "Age", Err: ErrMin}},
	},
	{
		in: User{
			ID:     "123456",                               // string `json:"id" validate:"len:6"`
			Age:    20,                                     // int      `validate:"min:18|max:50"`
			Email:  "batonhleba.com",                       // string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
			Role:   "admin",                                // UserRole `validate:"in:admin,stuff"`
			Phones: []string{"89125885738", "89125885736"}, // []string `validate:"len:11"`
		},
		expectedErr: ValidationErrors{{Field: "Email", Err: ErrRegexp}},
	},
	{
		in: User{
			ID:     "123456",                               // string `json:"id" validate:"len:6"`
			Age:    20,                                     // int      `validate:"min:18|max:50"`
			Email:  "baton@hleba.com",                      // string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
			Role:   "huylo",                                // UserRole `validate:"in:admin,stuff"`
			Phones: []string{"89125885738", "89125885736"}, // []string `validate:"len:11"`
		},
		expectedErr: ValidationErrors{{Field: "Role", Err: ErrIn}},
	},
	{
		in: User{
			ID:     "123456",                                // string `json:"id" validate:"len:6"`
			Age:    20,                                      // int      `validate:"min:18|max:50"`
			Email:  "baton@hleba.com",                       // string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
			Role:   "admin",                                 // UserRole `validate:"in:admin,stuff"`
			Phones: []string{"891825885738", "89125885736"}, // []string `validate:"len:11"`
		},
		expectedErr: ValidationErrors{{Field: "Phones", Err: ErrLen}},
	},
	{
		in: User{
			ID:     "12456",                             // string `json:"id" validate:"len:6"`
			Age:    2,                                   // int      `validate:"min:18|max:50"`
			Email:  "baton@hleba.com",                   // string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
			Role:   "CEO",                               // UserRole `validate:"in:admin,stuff"`
			Phones: []string{"89008985738", "89285736"}, // []string `validate:"len:11"`
		},
		expectedErr: ValidationErrors{
			{Field: "ID", Err: ErrLen},
			{Field: "Age", Err: ErrMin},
			{Field: "Role", Err: ErrIn},
			{Field: "Phones", Err: ErrLen},
		},
	},
	{
		in:          App{Version: "12345"},
		expectedErr: nil,
	},
	{
		in:          App{Version: "123456"},
		expectedErr: ValidationErrors{{Field: "Version", Err: ErrLen}},
	},
	{
		in:          App{Version: ""},
		expectedErr: ValidationErrors{{Field: "Version", Err: ErrLen}},
	},

	{
		in:          Token{},
		expectedErr: nil,
	},

	{
		in: Response{
			Code: 200, // int    `validate:"in:200,404,500"`
			Body: "",  // string `json:"omitempty"`
		},
		expectedErr: nil,
	},
	{
		in: Response{
			Code: 0,  // int    `validate:"in:200,404,500"`
			Body: "", // string `json:"omitempty"`
		},
		expectedErr: ValidationErrors{{Field: "Code", Err: ErrIn}},
	},

	{in: Malformed1{}, expectedErr: ErrTag},
	{in: Malformed2{}, expectedErr: ErrTag},
	{in: Malformed3{}, expectedErr: ErrTag},
	{in: Malformed4{}, expectedErr: ErrTag},

	{in: MultiFailed1{Field1: 8}, expectedErr: ValidationErrors{
		{Field: "Field1", Err: ErrMin},
		{Field: "Field1", Err: ErrMin},
		{Field: "Field1", Err: ErrMin},
		{Field: "Field1", Err: ErrMax},
		{Field: "Field1", Err: ErrIn},
	}},

	{
		in: StructTest1{
			Field1: 1, // err
			Field2: StructTest2{
				Field2_1: StructTest3{
					Field3_1: 1,   // ok
					Field3_2: "5", // err
				},
				Field2_2: 11, // err
			},
		},
		expectedErr: ValidationErrors{
			{Field: "Field1", Err: ErrMin},
			{Field: "Field3_2", Err: ErrIn},
			{Field: "Field2_2", Err: ErrMax},
		},
	},

	{in: StructTest4{}, expectedErr: ErrTag},
	{in: []string{}, expectedErr: ErrNotAStruct},
}

func TestValidate(t *testing.T) {
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)

			if tt.expectedErr == nil {
				require.Nil(t, err)
				return
			}

			var expErrs ValidationErrors
			if errors.As(tt.expectedErr, &expErrs) {
				var resultErrs ValidationErrors
				require.ErrorAs(t, err, &resultErrs)
				require.Equal(t, len(expErrs), len(resultErrs))

				for j := range expErrs {
					require.ErrorIs(t, expErrs[j].Err, resultErrs[j].Err)
					require.Equal(t, expErrs[j].Field, resultErrs[j].Field)
				}
			} else {
				require.ErrorIs(t, tt.expectedErr, err)
			}
		})
	}
}
