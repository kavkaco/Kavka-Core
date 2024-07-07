package vali

import "github.com/go-playground/validator/v10"

type Varror struct {
	Error            error
	ValidationErrors validator.ValidationErrors
}
