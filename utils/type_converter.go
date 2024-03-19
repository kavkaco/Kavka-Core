package utils

import (
	"go.mongodb.org/mongo-driver/bson"
)

func TypeConverter[R any](data any) (*R, error) {
	var result R

	b, err := bson.Marshal(&data)
	if err != nil {
		return nil, err
	}

	err = bson.Unmarshal(b, &result)
	if err != nil {
		return nil, err
	}

	return &result, err
}
