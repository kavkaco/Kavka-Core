package structs

import (
	"errors"
	"reflect"
	"strings"
)

var (
	ErrMustBePointer = errors.New("cannot assign to the item passed, item must be a pointer in order to assign")
	ErrTagNotFound   = errors.New("tag provided does not define a tag")
	ErrFieldNotFound = errors.New("field does not exist within the provided item")
)

func SetFieldByBSON(item interface{}, fieldName string, value interface{}) error {
	v := reflect.ValueOf(item).Elem()
	if !v.CanAddr() {
		return ErrMustBePointer
	}

	findBSONName := func(t reflect.StructTag) (string, error) {
		if jt, ok := t.Lookup("bson"); ok {
			return strings.Split(jt, ",")[0], nil
		}
		return "", ErrTagNotFound
	}

	fieldNames := map[string]int{}

	for i := 0; i < v.NumField(); i++ {
		typeField := v.Type().Field(i)
		tag := typeField.Tag
		bsonName, _ := findBSONName(tag)
		fieldNames[bsonName] = i
	}

	fieldNum, ok := fieldNames[fieldName]
	if !ok {
		return ErrFieldNotFound
	}

	fieldVal := v.Field(fieldNum)
	fieldVal.Set(reflect.ValueOf(value))

	return nil
}
