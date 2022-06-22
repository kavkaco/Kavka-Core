package controllers

import (
	"Tahagram/database"
	"Tahagram/httpstatus"
	"Tahagram/models"
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

type SigninBody struct {
	Email string `json:"email" xml:"email" form:"email" validate:"required,email"`
}

type User struct {
	Username string `json:"username"`
}

func SigninAction(c *fiber.Ctx) error {
	u := new(SigninBody)

	if err := c.BodyParser(u); err != nil {
		httpstatus.InternalServerError(c)
		return err
	}

	errors := models.ValidateStruct(u)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	filter := bson.D{{"email", u.Email}}
	result := database.UsersCollection.FindOne(context.TODO(), filter)

	fmt.Println(result)
	// ANCHOR

	return nil
}

func VerifyCodeAction(c *fiber.Ctx) error {
	return nil
}

func AuthenticationAction(c *fiber.Ctx) error {

	return nil
}

// sess, sessErr := session.SessionStore.Get(c)
// if sessErr != nil {
// 	// FIXME - fix auto error handlers
// 	c.Status(500).SendString("Internal Server Error")
// }

// fmt.Println(sess)
