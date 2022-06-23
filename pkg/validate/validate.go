package validate

import (
	"fmt"
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

	for i := 0; i < t.NumField(); i++ { // each field
		field := t.Field(i)
		allTags := field.Tag.Get(validateTagName)
		value := reflect.ValueOf(inter).Field(i).String()

		fmt.Printf("%s:%s\n", field.Name, value)

		tags := strings.Split(allTags, ";")
		for _, tag := range tags { // each tag
			// fmt.Printf("Field:%s\n", field.Name)
			// fmt.Printf("Value:%s\n", value)
			// fmt.Println("-----------")

			switch tag {
			case "required":
				if len(strings.TrimSpace(value)) == 0 {
					errors = append(errors, ValidationError{
						Message: "required",
						Field:   field.Name,
					})
				}
			case "email":
				if !mailRegex.MatchString(value) {
					errors = append(errors, ValidationError{
						Message: "not valid",
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
