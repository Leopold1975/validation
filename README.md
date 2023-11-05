# Validation

[![Checks status](https://github.com/Leopold1975/validation/actions/workflows/tests.yml/badge.svg)](https://github.com/Leopold1975/validation/actions/workflows/tests.yml)

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