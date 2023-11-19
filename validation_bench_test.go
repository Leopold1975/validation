package validation_test

import (
	"encoding/json"
	"log"
	"os"
	"testing"

	"github.com/Leopold1975/validation"
	"github.com/asaskevich/govalidator"
)

func BenchmarkValidate(b *testing.B) {

	f, err := os.OpenFile("tests.txt", os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	dec := json.NewDecoder(f)

	benchCases := []TestCase{}
	if err := dec.Decode(&benchCases); err != nil {
		log.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		for _, tt := range benchCases {
			err := validation.Validate(tt)
			_ = err
		}
	}
}

func BenchmarkGoValidate(b *testing.B) {

	f, err := os.OpenFile("tests.txt", os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	dec := json.NewDecoder(f)

	benchCases := []TestCase{}
	if err := dec.Decode(&benchCases); err != nil {
		log.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		for _, tt := range benchCases {
			_, err := govalidator.ValidateStruct(tt)
			_ = err
		}
	}
}
