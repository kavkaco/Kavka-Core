package models

import (
	"reflect"
	"strings"
)

type ValidationError struct {
	Message string
	Field   string
}

// var mailRegex = regexp.MustCompile(`\A[\w+\-.]+@[a-z\d\-]+(\.[a-z]+)*\.[a-z]+\z`)

func ValidateStruct(inter interface{}) []ValidationError {
	const validateTagName = "validate"
	var errors []ValidationError

	t := reflect.TypeOf(inter)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		allTags := field.Tag.Get(validateTagName)
		values := reflect.ValueOf(inter)

		tags := strings.Split(allTags, ";")
		for tagIndex, tag := range tags {
			value := values.Field(tagIndex).String()

			switch tag {
			case "required":
				if len(strings.TrimSpace(value)) == 0 {
					errors = append(errors, ValidationError{
						Message: "required",
						Field:   field.Name,
					})
				}
			default:
				panic("ValidateStruct:->Invalid Tag!-> " + tag)
			}
		}
	}

	return errors
}
