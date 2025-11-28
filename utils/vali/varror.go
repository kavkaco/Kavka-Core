package vali

import "github.com/go-playground/validator/v10"

type ValiErr struct {
	Error            error
	ValidationErrors validator.ValidationErrors
}
