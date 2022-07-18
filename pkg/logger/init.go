package logger

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

var ErrorLogger *log.Logger

func InitLogger(app *fiber.App) {
	wd, _ := os.Getwd()
	logsDir := wd + "/logs"

	requestsFile, requestsFileErr := os.OpenFile(fmt.Sprintf("%s/requests.log", logsDir), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if requestsFileErr != nil {
		fmt.Println(requestsFileErr)
		os.Exit(1)
	}

	app.Use(
		logger.New(
			logger.Config{
				Format:   "Pid: ${pid}\nStatus: ${status}\nMethod: ${method}\nPath: ${path}\nTime: ${time} \n\n",
				TimeZone: "Asia/Iran",
				Output:   requestsFile,
			},
		),
	)

	errorsFile, errorsFileErr := os.OpenFile(fmt.Sprintf("%s/errors.log", logsDir), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if errorsFileErr != nil {
		fmt.Println(errorsFileErr)
		os.Exit(1)
	}

	ErrorLogger = log.New(errorsFile, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}
