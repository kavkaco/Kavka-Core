package controllers

import (
	"Tahagram/database"
	"Tahagram/httpstatus"
	"Tahagram/models"
	"Tahagram/pkg/auth"
	"Tahagram/pkg/validate"
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SigninBody struct {
	Email string `json:"email" xml:"email" form:"email" validate:"required;email"`
}

type User struct {
	Username string `json:"username"`
}

func SigninAction(c *fiber.Ctx) error {
	var body SigninBody

	if err := c.BodyParser(&body); err != nil {
		httpstatus.InternalServerError(c)
		return err
	}

	errors := validate.ValidateStruct(body)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	var user *models.User

	filter := bson.D{
		primitive.E{Key: "email", Value: body.Email},
	}
	database.UsersCollection.FindOne(context.TODO(), filter).Decode(&user)

	if user != nil {
		if !user.VerificCodeLimitDate.IsZero() && auth.UserLimited(&user.VerificCodeLimitDate) {
			c.Status(200).JSON(fiber.Map{
				"Message":  "limited",
				"LimitEnd": user.VerificCodeLimitDate.String(),
			})
		} else {
			newData := bson.D{
				primitive.E{Key: "$set", Value: bson.D{
					primitive.E{Key: "verific_code", Value: uint(auth.MakeVerificCode())},
					primitive.E{Key: "verific_code_expire", Value: auth.MakeVerificCodeExpire(time.Now())},
					primitive.E{Key: "verific_code_limit_date", Value: time.Time{}},
				}},
			}

			database.UsersCollection.FindOneAndUpdate(context.TODO(), bson.D{
				primitive.E{Key: "email", Value: body.Email},
			}, newData)

			err := sendSigninEmail()
			if err != nil {
				httpstatus.InternalServerError(c)
			} else {
				c.Status(200).JSON(fiber.Map{
					"Message": "verific code sent",
				})
			}

		}
	} else {
		sendEmailErr := sendSigninEmail()
		if sendEmailErr != nil {
			httpstatus.InternalServerError(c)
		} else {
			_, insertErr := database.UsersCollection.InsertOne(context.TODO(), bson.D{
				primitive.E{Key: "email", Value: body.Email},
				primitive.E{Key: "username", Value: auth.GetEmailWithoutAt(body.Email)},
				primitive.E{Key: "verific_code", Value: uint(auth.MakeVerificCode())},
				primitive.E{Key: "verific_try_count", Value: 0},
			})
			if insertErr != nil {
				httpstatus.InternalServerError(c)
			} else {
				c.Status(200).JSON(fiber.Map{
					"Message": "verific code sent",
				})
			}
		}
	}

	return nil
}

func VerifyCodeAction(c *fiber.Ctx) error {
	return nil
}

func AuthenticationAction(c *fiber.Ctx) error {

	return nil
}

func sendSigninEmail() error {
	fmt.Println("Email sent")
	return nil
}

// sess, sessErr := session.SessionStore.Get(c)
// if sessErr != nil {
// 	// FIXME - fix auto error handlers
// 	c.Status(500).SendString("Internal Server Error")
// }

// fmt.Println(sess)
