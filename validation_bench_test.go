package validation_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/Leopold1975/validation"
	"github.com/asaskevich/govalidator"
	"github.com/stretchr/testify/require"
)

func BenchmarkValidate(b *testing.B) {
	f, err := os.OpenFile("tests.txt", os.O_CREATE|os.O_RDWR, 0o755)
	require.Nil(b, err)

	defer f.Close()
	dec := json.NewDecoder(f)

	benchCases := []TestCase{}
	err = dec.Decode(&benchCases)
	require.Nil(b, err)

	for i := 0; i < b.N; i++ {
		for _, tt := range benchCases {
			err := validation.Validate(tt)
			_ = err
		}
	}
}

func BenchmarkGoValidate(b *testing.B) {
	f, err := os.OpenFile("tests.txt", os.O_CREATE|os.O_RDWR, 0o755)
	require.Nil(b, err)

	defer f.Close()
	dec := json.NewDecoder(f)

	benchCases := []TestCase{}
	err = dec.Decode(&benchCases)
	require.Nil(b, err)

	for i := 0; i < b.N; i++ {
		for _, tt := range benchCases {
			_, err := govalidator.ValidateStruct(tt)
			_ = err
		}
	}
}
