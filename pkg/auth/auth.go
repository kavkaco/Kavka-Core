package auth

import (
	"Tahagram/configs"
	"Tahagram/database"
	"Tahagram/httpstatus"
	"Tahagram/logs"
	"Tahagram/models"
	"Tahagram/session"
	"context"
	"math/rand"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func MakeVerificCode() int {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)

	max := 999999
	min := 111111

	return min + r.Intn(max-min)
}

func MakeVerificCodeExpire() int64 {
	return time.Now().Add(configs.VerificCodeExpire).Unix()
}

func MakeVerificLimitDate() int64 {
	return time.Now().Add(configs.VerificLimitDate).Unix()
}

func IsUserLimited(limitDate int64) bool {
	now := time.Now().Unix()
	return !(now < limitDate)
}

func IsVerificCodeExpired(expire int64) bool {
	now := time.Now().Unix()
	return !(now < expire)
}

func GetEmailWithoutAt(email string) string {
	return email[:strings.IndexByte(email, '@')]
}

func AuthenticateUser(c *fiber.Ctx) (bool, *models.User) {
	sess, sessErr := session.SessionStore.Get(c)
	if sessErr != nil {
		logs.ErrorLogger.Println(sessErr)
		httpstatus.InternalServerError(c)
	}

	userId := sess.Get("static_id")

	if userId != nil && userId.(int) > 0 {
		var user *models.User

		database.UsersCollection.FindOne(context.TODO(), bson.D{
			primitive.E{
				Key:   "static_id",
				Value: userId.(int),
			},
		}).Decode(&user)

		if user != nil {
			return true, user
		} else {
			return false, nil
		}
	} else {
		return false, nil
	}
}
