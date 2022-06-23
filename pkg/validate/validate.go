package validate

import (
	"reflect"
	"regexp"
	"strings"
)

type ValidationError struct {
	Message string
	Field   string
}

var mailRegex = regexp.MustCompile(`\A[\w+\-.]+@[a-z\d\-]+(\.[a-z]+)*\.[a-z]+\z`)

func ValidateStruct(inter interface{}) []ValidationError {
	const validateTagName = "validate"
	var errors []ValidationError

	t := reflect.TypeOf(inter)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		allTags := field.Tag.Get(validateTagName)
		tags := strings.Split(allTags, ";")
		value := reflect.ValueOf(inter).Field(i).String()
		for _, tag := range tags {
			switch tag {
			case "required":
				if len(strings.TrimSpace(value)) == 0 {
					errors = append(errors, ValidationError{
						Message: "required",
						Field:   strings.ToLower(field.Name),
					})
				}
			case "email":
				if !mailRegex.MatchString(value) {
					errors = append(errors, ValidationError{
						Message: "not valid",
						Field:   strings.ToLower(field.Name),
					})
				}
			}
		}
	}

	return errors
}
