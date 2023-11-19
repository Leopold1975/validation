# Validation

[![Tests status](https://github.com/Leopold1975/validation/actions/workflows/tests.yml/badge.svg)](https://github.com/Leopold1975/validation/actions/workflows/tests.yml)

The Validation project is a Go library for structure validation.
Validation allows you to define validation rules for structure fields and check if the values stick to the rules.

## Install

Install the validation library using 
``` 
go get github.com/Leopold1975/validation
```

## Usage

An example of using the validation library:

```go
package main

import (
    "fmt"
    "github.com/Leopold1975/validation"
)

type User struct {
    Name  string `validate:"in:Arthur,Alex"`
    Email string `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
    Age   int    `validate:"min:18"`
    Role string  `validate:"in:admin,user"`
}

func main() {
    user := User{
        Name:  "John Doe",
        Email: "johndoe@example.com",
        Age:   25,
        Role: "admin",
    }

    err := validation.Validate(user)
    if err != nil {
        fmt.Printf("Validation error: %v\n", err)
    } else {
        fmt.Println("Validation successful!")
    }
}

```

## Compare to govalidator

Benchmarks for `Validate` and `ValidateStruct` from [govalidator](https://github.com/asaskevich/govalidator#validatestruct-2).

```bash
goos: linux
goarch: amd64
pkg: github.com/Leopold1975/validation
cpu: 11th Gen Intel(R) Core(TM) i5-11400H @ 2.70GHz
              │ leopold_mem.out │ govalidator_mem.out │
              │     sec/op      │    sec/op     vs base   │
Validate-12        185.9µ ± 15%
GoValidate-12                      1.144m ± 3%
geomean            185.9µ          1.144m       ? ¹ ²
¹ benchmark set differs from baseline; geomeans may not be comparable
² ratios must be >0 to compute geomean

              │ leopold_mem.out │ govalidator_mem.out │
              │      B/op       │     B/op      vs base   │
Validate-12        83.24Ki ± 0%
GoValidate-12                     343.0Ki ± 0%
geomean            83.24Ki        343.0Ki       ? ¹ ²
¹ benchmark set differs from baseline; geomeans may not be comparable
² ratios must be >0 to compute geomean

              │ leopold_mem.out │ govalidator_mem.out │
              │    allocs/op    │  allocs/op    vs base   │
Validate-12         3.042k ± 0%
GoValidate-12                      15.82k ± 0%
geomean             3.042k         15.82k       ? ¹ ²
```

## Contribution

If you have any suggestions or fixes, feel free to contribute. The project is open for pull requests and discussions.
