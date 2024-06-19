package repository_mongo

import (
	"context"
	"errors"
	"time"

	"github.com/kavkaco/Kavka-Core/database"
	"github.com/kavkaco/Kavka-Core/internal/model"
	"github.com/kavkaco/Kavka-Core/internal/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type authRepository struct {
	authCollection *mongo.Collection
}

func NewAuthMongoRepository(db *mongo.Database) repository.AuthRepository {
	return &authRepository{db.Collection(database.AuthCollection)}
}

func (a *authRepository) IncrementFailedLoginAttempts(ctx context.Context, userID string) error {
	filter := bson.M{"user_id": userID}
	update := bson.M{"$inc": bson.M{"failed_login_attempts": 1}}

	_, err := a.authCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (a *authRepository) ClearFailedLoginAttempts(ctx context.Context, userID model.UserID) error {
	filter := bson.M{"user_id": userID}
	update := bson.M{"$set": bson.M{"failed_login_attempts": 0}}

	_, err := a.authCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (a *authRepository) LockAccount(ctx context.Context, userID model.UserID, lockDuration time.Duration) error {
	now := time.Now()
	now = now.Add(lockDuration)

	filter := bson.M{"user_id": userID}
	update := bson.M{"$set": bson.M{"account_locked_until": now.Unix()}}

	_, err := a.authCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (a *authRepository) UnlockAccount(ctx context.Context, userID model.UserID) error {
	filter := bson.M{"user_id": userID}
	update := bson.M{"$set": bson.M{"account_locked_until": 0}}

	_, err := a.authCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (a *authRepository) VerifyEmail(ctx context.Context, userID model.UserID) error {
	filter := bson.M{"user_id": userID}
	update := bson.M{"$set": bson.M{"email_verified": true}}

	_, err := a.authCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return err
}

func (a *authRepository) Create(ctx context.Context, authModel *model.Auth) (*model.Auth, error) {
	_, err := a.authCollection.InsertOne(ctx, authModel)
	if err != nil {
		return nil, err
	}

	return authModel, nil
}

func (a *authRepository) GetUserAuth(ctx context.Context, userID model.UserID) (*model.Auth, error) {
	result := a.authCollection.FindOne(ctx, bson.M{"user_id": userID})
	if errors.Is(result.Err(), mongo.ErrNoDocuments) {
		return nil, repository.ErrNotFound
	} else if result.Err() != nil {
		return nil, result.Err()
	}

	var authModel model.Auth
	err := result.Decode(&authModel)
	if err != nil {
		return nil, err
	}

	return &authModel, nil
}

func (a *authRepository) ChangePassword(ctx context.Context, userID model.UserID, passwordHash string) error {
	result := a.authCollection.FindOneAndUpdate(ctx, bson.M{"user_id": userID}, bson.M{"$set": bson.M{"password_hash": passwordHash}})
	if result.Err() != nil {
		return result.Err()
	}

	return nil
}

func (repo *authRepository) DeleteByID(ctx context.Context, userID model.UserID) error {
	filter := bson.M{"user_id": userID}

	result, err := repo.authCollection.DeleteOne(ctx, filter)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return repository.ErrNotFound
	} else if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return repository.ErrNotDeleted
	}

	return nil
}
