package main

import (
	"Tahagram/models"
	"testing"
)

func TestStructToBson(t *testing.T) {
	u := models.User{
		Name:     "Taha",
		Username: "tahadostifam",
	}

	bsonDoc, err := models.StructToBSON(u)
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log(bsonDoc)
	}
}
