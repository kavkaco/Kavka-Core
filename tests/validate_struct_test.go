package main

import (
	"Kavka/pkg/validate"
	"os"
	"testing"

	"github.com/olekukonko/tablewriter"
)

type User struct {
	Name  string `validate:"required"`
	Email string `validate:"email"`
}

func TestValidateStruct(t *testing.T) {
	u := User{
		Name:  "Taha",
		Email: "taha@mail.com",
	}

	errors := validate.ValidateStruct(u)
	if errors != nil {
		writeValidationErrorsToTable(errors)
	} else {
		t.Log("Valid Struct!")
	}
}

func writeValidationErrorsToTable(errors []validate.ValidationError) {
	var data = [][]string{}
	for _, v := range errors {
		data = append(data, []string{v.Field, v.Message})
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Field", "Message"})
	for _, v := range data {
		table.Append(v)
	}
	table.Render()
}
