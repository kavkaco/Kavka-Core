package httpstatus

import (
	"Nexus/app/models"

	"github.com/fatih/structs"
	"github.com/gofiber/fiber/v2"
)

func InternalServerError(c *fiber.Ctx) {
	c.Status(500).JSON(fiber.Map{
		"message": "An error occurred on the server side",
	})
}

func Unauthorized(c *fiber.Ctx) {
	c.Status(401).JSON(fiber.Map{
		"message": "Unauthorized",
	})
}

func ResponseUserData(c *fiber.Ctx, user *models.User) {
	var dataToSend map[string]interface{} = structs.Map(user)
	delete(dataToSend, "VerificCode")
	delete(dataToSend, "VerificTryCount")
	delete(dataToSend, "VerificCodeExpire")
	delete(dataToSend, "VerificLimitDate")
	delete(dataToSend, "FirstLoginCompleted")

	c.Status(200).JSON(dataToSend)
}
