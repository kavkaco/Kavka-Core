package logs

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func InitLogger(app *fiber.App) {
	wd, _ := os.Getwd()
	file, err := os.OpenFile(wd+"/logs/info.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	app.Use(
		logger.New(
			logger.Config{
				Format:   "Pid: ${pid}\nStatus: ${status}\nMethod: ${method}\nPath: ${path}\nTime: ${time} \n\n",
				TimeZone: "Asia/Iran",
				Output:   file,
			},
		),
	)
}
