package controllers

import (
	"Kavka/app/database"
	"Kavka/app/httpstatus"
	"Kavka/app/models"
	"Kavka/app/session"
	"Kavka/app/smtp"
	"Kavka/internal/configs"
	"Kavka/pkg/auth"
	"Kavka/pkg/logger"
	"Kavka/pkg/validate"
	"context"
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
type AuthBody struct {
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

	var user *models.User = models.FindUserByEmail(body.Email)

	if user != nil {
		if !auth.IsUserLimited(user.VerificLimitDate) {
			c.Status(200).JSON(fiber.Map{
				"Message":  "limited",
				"LimitEnd": time.Unix(user.VerificLimitDate, 0).String(),
			})
		} else {
			newData := bson.D{
				primitive.E{Key: "$set", Value: bson.D{
					primitive.E{Key: "verific_code", Value: uint(auth.MakeVerificCode())},
					primitive.E{Key: "verific_code_expire", Value: auth.MakeVerificCodeExpire()},
					primitive.E{Key: "verific_limit_date", Value: nil},
				}},
			}

			database.UsersCollection.FindOneAndUpdate(context.TODO(), bson.D{
				primitive.E{Key: "email", Value: body.Email},
			}, newData)

			sendEmailErr := smtp.SendSigninEmail()
			if sendEmailErr != nil {
				logger.ErrorLogger.Println(sendEmailErr)
				httpstatus.InternalServerError(c)
			} else {
				c.Status(200).JSON(fiber.Map{
					"Message": "verific code sent",
				})
			}

		}
	} else {
		sendEmailErr := smtp.SendSigninEmail()
		if sendEmailErr != nil {
			logger.ErrorLogger.Println(sendEmailErr)
			httpstatus.InternalServerError(c)
		} else {
			_, insertErr := database.UsersCollection.InsertOne(context.TODO(), bson.D{
				primitive.E{Key: "static_id", Value: models.MakeUserStaticID()},
				primitive.E{Key: "email", Value: body.Email},
				primitive.E{Key: "username", Value: auth.GetEmailWithoutAt(body.Email)},
				primitive.E{Key: "verific_code", Value: uint(auth.MakeVerificCode())},
				primitive.E{Key: "verific_try_count", Value: 0},
				primitive.E{Key: "verific_code_expire", Value: auth.MakeVerificCodeExpire()},
			})
			if insertErr != nil {
				logger.ErrorLogger.Println(insertErr)
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
		if !auth.IsUserLimited(user.VerificLimitDate) {
			c.Status(200).JSON(fiber.Map{
				"Message":  "limited",
				"LimitEnd": time.Unix(user.VerificLimitDate, 0).String(),
			})
		} else {
			if user.VerificTryCount >= configs.MaxVerificTryCount {
				limit := auth.MakeVerificLimitDate()
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
									Value: limit,
								},
								primitive.E{
									Key:   "verific_try_count",
									Value: 0,
								},
							},
						},
					},
				)

				c.Status(200).JSON(fiber.Map{
					"Message":  "limited",
					"LimitEnd": time.Unix(limit, 0),
				})
			} else {
				userVerificCode, userVerificCodeErr := strconv.Atoi(body.VerificCode)
				if userVerificCodeErr != nil {
					httpstatus.Unauthorized(c)
				} else {
					if user.VerificCode == userVerificCode {
						if !auth.IsVerificCodeExpired(user.VerificCodeExpire) {
							sess, sessErr := session.SessionStore.Get(c)
							if sessErr != nil {
								user.IncreaseTryCount()
								logger.ErrorLogger.Println(sessErr)
								httpstatus.InternalServerError(c)
							}

							sess.Set("static_id", user.StaticID)

							database.UsersCollection.FindOneAndUpdate(context.TODO(),
								bson.D{
									primitive.E{Key: "email", Value: body.Email},
								},
								bson.D{
									primitive.E{
										Key: "$set",
										Value: bson.D{
											primitive.E{Key: "verific_code", Value: nil},
											primitive.E{Key: "verific_try_count", Value: nil},
											primitive.E{Key: "verific_code_expire", Value: nil},
											primitive.E{Key: "verific_limit_date", Value: nil},
											primitive.E{Key: "first_login_completed", Value: true},
										},
									},
								},
							)

							defer sess.Save()

							c.Status(200).JSON(fiber.Map{
								"Message": "success signin",
							})
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
		}
	} else {
		httpstatus.Unauthorized(c)
	}

	return nil
}

func AuthenticationAction(c *fiber.Ctx) error {
	authorized, user := auth.AuthenticateUser(c)

	if authorized {
		httpstatus.ResponseUserData(c, user)
	} else {
		httpstatus.Unauthorized(c)
	}

	return nil
}

func LogoutAction(c *fiber.Ctx) error {
	sess, sessErr := session.SessionStore.Get(c)
	if sessErr != nil {
		logger.ErrorLogger.Println(sessErr)
		httpstatus.InternalServerError(c)
	}

	err := sess.Destroy()

	if err != nil {
		httpstatus.InternalServerError(c)
	} else {
		c.Status(200).JSON(fiber.Map{
			"Message": "logout",
		})
	}

	return nil
}
