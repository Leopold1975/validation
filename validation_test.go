package validation_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/Leopold1975/validation"
	"github.com/stretchr/testify/require"
)

type UserRole string

type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff|len:5"`
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
	TestStruct struct {
		TestStringLen   string   `validate:"len:4"`
		TestStringSlice []string `validate:"len:2"`
		TestInStrings   string   `validate:"in:foo,bar"`
		TestInInt       int      `validate:"in:16,24"`
		TestInSliceStr  []string `validate:"in:foo,bar"`
		TestInSliceInt  []int    `validate:"in:16,8"`
		TestRegexp      string   `validate:"regexp:^\\d+$"`
		TestRegexpSlice []string `validate:"regexp:^[A-Z]+[a-z]+$"`
		TestIntMin      int      `validate:"min:8"`
		TestIntMax      int      `validate:"max:8"`
		TestIntMinSlice []int    `validate:"min:5"`
		TestIntMaxSlice []int    `validate:"max:5"`
	}
	NestedStruct struct {
		Resp Response
		Name string `validate:"len:6"`
	}
)

var testCases = []struct {
	in          interface{}
	expectedErr error
}{
	{
		TestStruct{
			TestStringLen:   "good",
			TestStringSlice: []string{"aa", "bb", "cc"},
			TestInStrings:   "bar",
			TestInInt:       16,
			TestInSliceStr:  []string{"bar", "foo"},
			TestInSliceInt:  []int{8, 16, 8},
			TestRegexp:      "1234",
			TestRegexpSlice: []string{"Joe", "Alex", "Max"},
			TestIntMin:      9,
			TestIntMax:      7,
			TestIntMinSlice: []int{5, 6, 7},
			TestIntMaxSlice: []int{1, 2, 3, 4},
		},
		nil,
	},
	{
		TestStruct{
			TestStringLen:   "wrong",
			TestStringSlice: []string{"foo", "bar", "ccc"},
			TestInStrings:   "foe",
			TestInInt:       17,
			TestInSliceStr:  []string{"abcbar", "cde", "foefoor"},
			TestInSliceInt:  []int{8, 16, 9},
			TestRegexp:      "1a2b3cd5",
			TestRegexpSlice: []string{"Joe", "Alex", "max"},
			TestIntMin:      7,
			TestIntMax:      9,
			TestIntMinSlice: []int{1, 2, 3, 4},
			TestIntMaxSlice: []int{6, 7, 8},
		},
		validation.Errors{
			validation.Error{
				Field: "TestStringLen",
				Err:   validation.ErrLenNotEqual,
			},
			validation.Error{
				Field: "TestStringSlice",
				Err:   validation.ErrLenNotEqual,
			},
			validation.Error{
				Field: "TestInStrings",
				Err:   validation.ErrOutOfSet,
			},
			validation.Error{
				Field: "TestInInt",
				Err:   validation.ErrOutOfSet,
			},
			validation.Error{
				Field: "TestInSliceStr",
				Err:   validation.ErrOutOfSet,
			},
			validation.Error{
				Field: "TestInSliceInt",
				Err:   validation.ErrOutOfSet,
			},
			validation.Error{
				Field: "TestRegexp",
				Err:   validation.ErrNotRegexp,
			},
			validation.Error{
				Field: "TestRegexpSlice",
				Err:   validation.ErrNotRegexp,
			},
			validation.Error{
				Field: "TestIntMin",
				Err:   validation.ErrValueOutOfBondary,
			},
			validation.Error{
				Field: "TestIntMax",
				Err:   validation.ErrValueOutOfBondary,
			},
			validation.Error{
				Field: "TestIntMinSlice",
				Err:   validation.ErrValueOutOfBondary,
			},
			validation.Error{
				Field: "TestIntMaxSlice",
				Err:   validation.ErrValueOutOfBondary,
			},
		},
	},
	{
		User{
			ID:     "1",
			Name:   "John",
			Age:    51,
			Email:  "assfag@asfa",
			Role:   "aaaaa",
			Phones: []string{"79998886655", "+114442223311", "1234"},
			meta:   json.RawMessage{},
		},
		validation.Errors{
			validation.Error{
				Field: "ID",
				Err:   validation.ErrLenNotEqual,
			},
			validation.Error{
				Field: "Age",
				Err:   validation.ErrValueOutOfBondary,
			},
			validation.Error{
				Field: "Email",
				Err:   validation.ErrNotRegexp,
			},
			validation.Error{
				Field: "Role",
				Err:   validation.ErrOutOfSet,
			},
			validation.Error{
				Field: "Phones",
				Err:   validation.ErrLenNotEqual,
			},
		},
	},
	{
		User{
			ID:     "GoodUserGoodUserGoodUserGoodUserGood",
			Age:    19,
			Email:  "mefisto123@secretdoami.com",
			Role:   "admin",
			Phones: []string{"79998886655", "84442223311", "18664443322"},
			meta:   json.RawMessage{1, 2, 3},
		},
		nil,
	},
	{
		User{
			ID:    "PartiallyGooduser",
			Age:   50,
			Email: "mefisto123@secretdoami.com",
			Role:  "not_admin",
		},
		validation.Errors{
			validation.Error{
				Field: "ID",
				Err:   validation.ErrLenNotEqual,
			},
			validation.Error{
				Field: "Role",
				Err:   validation.ErrOutOfSet,
			},
		},
	},
	{
		5,
		validation.ErrUnsupportedType,
	},
	{
		App{
			Version: "1.3.5",
		},
		nil,
	},
	{
		App{
			Version: "1",
		},
		validation.Errors{
			validation.Error{
				Field: "Version",
				Err:   validation.ErrLenNotEqual,
			},
		},
	},
	{
		Token{
			Header:  []byte{1, 2, 3},
			Payload: []byte{},
		},
		nil,
	},
	{
		Response{
			Code: 200,
			Body: "body",
		},
		nil,
	},
	{
		Response{
			Code: 201,
			Body: "body",
		},
		validation.Errors{
			validation.Error{
				Field: "Code",
				Err:   validation.ErrOutOfSet,
			},
		},
	},
	{
		NestedStruct{
			Resp: Response{
				Code: 200,
			},
			Name: "Struct",
		},
		nil,
	},
}

func TestValidate(t *testing.T) {
	for i, tt := range testCases {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			// t.Parallel()

			err := validation.Validate(tt.in)

			require.Equal(t, tt.expectedErr, err)
		})
	}
}

var testsErr = []struct {
	in          interface{}
	expectedErr error
}{
	{
		User{
			ID:    "PartiallyGooduser",
			Age:   50,
			Email: "mefisto123@secretdoami.com",
			Role:  "not_admin",
		},
		errors.New("field: ID, error: len doesn't satisfy the rule field: Role, error: value doesn't belong to the set "),
	},
	{
		struct {
			Ints []int8 `validate:"len:8"`
		}{
			[]int8{1, 2, 3},
		},
		validation.ErrUnsupportedType,
	},
	{
		struct {
			OneInt int `validate:"len:8"`
		}{
			1,
		},
		validation.ErrUnsupportedType,
	},
	{
		struct {
			NotGoodLen string `validate:"len:8a"`
		}{
			"abcd",
		},
		validation.ErrWrongValue,
	},
	{
		struct {
			NotGoodRegexp string `validate:"regexp:\\c+asavx"`
		}{
			"abcd",
		},
		validation.ErrWrongValue,
	},
	{
		struct {
			NotGoodRegexp []int `validate:"regexp:\\d+"`
		}{
			[]int{1, 2, 3},
		},
		validation.ErrUnsupportedType,
	},
	{
		struct {
			NotGoodRegexp int8 `validate:"regexp:\\d+"`
		}{
			3,
		},
		validation.ErrUnsupportedType,
	},
	{
		struct {
			NotGoodInt int `validate:"in:1d"`
		}{
			3,
		},
		validation.ErrWrongValue,
	},
	{
		struct {
			WrongTypeForIn int8 `validate:"in:1,2"`
		}{
			3,
		},
		validation.ErrUnsupportedType,
	},
	{
		struct {
			WrongTypeForMax int8 `validate:"max:1,2"`
		}{
			3,
		},
		validation.ErrUnsupportedType,
	},
	{
		struct {
			WrongSyntaxForMax int `validate:"max:1,2a"`
		}{
			3,
		},
		validation.ErrWrongValue,
	},
}

func TestErr(t *testing.T) {
	for i, tt := range testsErr {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			err := validation.Validate(tt.in)
			require.Equal(t, tt.expectedErr.Error(), err.Error())
		})
	}
}
