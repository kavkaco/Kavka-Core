package controllers

import (
	"Tahagram/configs"
	"Tahagram/database"
	"Tahagram/httpstatus"
	"Tahagram/models"
	"Tahagram/pkg/auth"
	"Tahagram/pkg/validate"
	"Tahagram/session"
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SigninBody struct {
	Email string `json:"email" xml:"email" form:"email" validate:"required;email"`
}

type VerifyCodeBody struct {
	Email       string `json:"email" xml:"email" form:"email" validate:"required;email"`
	VerificCode string `json:"code" xml:"code" form:"code" validate:"required"`
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

	var user *models.User = models.FindUserByEmail(body.Email)

	if user != nil {
		if !user.VerificLimitDate.IsZero() && auth.IsUserLimited(&user.VerificLimitDate) {
			c.Status(200).JSON(fiber.Map{
				"Message":  "limited",
				"LimitEnd": user.VerificLimitDate.String(),
			})
		} else {
			newData := bson.D{
				primitive.E{Key: "$set", Value: bson.D{
					primitive.E{Key: "verific_code", Value: uint(auth.MakeVerificCode())},
					primitive.E{Key: "verific_code_expire", Value: auth.MakeVerificCodeExpire()},
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
				primitive.E{Key: "static_id", Value: models.MakeUserStaticID()},
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
	var body VerifyCodeBody

	if err := c.BodyParser(&body); err != nil {
		httpstatus.InternalServerError(c)
		return err
	}

	errors := validate.ValidateStruct(body)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	var user *models.User = models.FindUserByEmail(body.Email)
	if user != nil {
		if user.VerificTryCount >= configs.MaxVerificTryCount {
			database.UsersCollection.UpdateOne(context.TODO(),
				bson.D{
					primitive.E{Key: "email", Value: body.Email},
				},
				bson.D{
					primitive.E{
						Key: "$set",
						Value: bson.D{
							primitive.E{
								Key:   "verific_limit_date",
								Value: auth.MakeVerificLimitDate(),
							},
						},
					},
				},
			)
			c.Status(503).JSON(fiber.Map{
				"Message": "maximum verific code try count",
			})
		} else {
			userVerificCode, userVerificCodeErr := strconv.Atoi(body.VerificCode)
			if userVerificCodeErr != nil {
				httpstatus.Unauthorized(c)
			} else {
				if user.VerificCode == userVerificCode {

					fmt.Println(time.Now())
					fmt.Println(user.VerificCodeExpire)
					fmt.Printf("Expired: %v\n", auth.IsVerificCodeExpired(user.VerificCodeExpire))

					if !auth.IsVerificCodeExpired(user.VerificCodeExpire) {
						sess, sessErr := session.SessionStore.Get(c)
						if sessErr != nil {
							user.IncreaseTryCount()
							httpstatus.InternalServerError(c)
						}

						sess.Set("user_id", user.StaticID)

						c.Status(200).JSON(fiber.Map{
							"Message": "success signin",
						})

						database.UsersCollection.FindOneAndUpdate(context.TODO(),
							bson.D{
								primitive.E{Key: "email", Value: body.Email},
							},
							bson.D{
								primitive.E{Key: "verific_code", Value: nil},
								primitive.E{Key: "verific_try_count", Value: nil},
								primitive.E{Key: "verific_code_expire", Value: nil},
								primitive.E{Key: "first_login_completed", Value: true},
							},
						)
					} else {
						user.IncreaseTryCount()
						c.Status(401).JSON(fiber.Map{
							"Message": "verific code expired",
						})
					}
				} else {
					user.IncreaseTryCount()
					httpstatus.Unauthorized(c)
				}
			}
		}
	} else {
		httpstatus.Unauthorized(c)
	}

	return nil
}

func AuthenticationAction(c *fiber.Ctx) error {

	return nil
}

func sendSigninEmail() error {
	fmt.Println("Email sent")
	return nil
}
