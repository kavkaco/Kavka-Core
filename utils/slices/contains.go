package slices

import "go.mongodb.org/mongo-driver/bson/primitive"

func ContainsObjectID(s []primitive.ObjectID, v primitive.ObjectID) bool {
	r := false
	for _, i := range s {
		if i.Hex() == v.Hex() {
			r = true
			break
		}
	}
	return r
}
