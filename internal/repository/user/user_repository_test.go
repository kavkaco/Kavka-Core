package repository

import (
	"testing"

	"github.com/kavkaco/Kavka-Core/internal/domain/user"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestUserRepository(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("test create", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateSuccessResponse())
		userRepo := NewUserRepository(mt.DB)

		model := user.NewUser("1234")
		model.Name = "John"
		model.LastName = "Doe"
		savedModel, err := userRepo.Create(model)

		assert.NoError(t, err)
		assert.Equal(t, savedModel.Name, model.Name)
	})

	mt.Run("test find by username", func(mt *mtest.T) {
		userRepo := NewUserRepository(mt.DB)

		expectedDoc := bson.D{
			{Key: "id", Value: primitive.NewObjectID()},
			{Key: "name", Value: "John"},
			{Key: "last_name", Value: "Doe"},
			{Key: "phone", Value: "1234"},
			{Key: "username", Value: "john_doe"},
		}

		mt.AddMockResponses(mtest.CreateCursorResponse(1, "myDB.users", mtest.FirstBatch, expectedDoc))

		model, err := userRepo.FindByUsername("john_doe")
		assert.NoError(t, err)

		assert.Equal(t, model.StaticID, expectedDoc.Map()["id"])
		assert.Equal(t, model.Name, expectedDoc.Map()["name"])
		assert.Equal(t, model.LastName, expectedDoc.Map()["last_name"])
		assert.Equal(t, model.Phone, expectedDoc.Map()["phone"])
	})
}
