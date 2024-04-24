package logs

import (
	"github.com/kavkaco/Kavka-Core/config"
	"go.uber.org/zap"
)

func InitZapLogger() *zap.Logger {
	var loggerConfig zap.Config

	if config.CurrentEnv == config.Development {
		loggerConfig = zap.NewProductionConfig()
		loggerConfig.OutputPaths = []string{config.ProjectRootPath + "/logs/logs.development.json"}
	} else {
		loggerConfig = zap.NewProductionConfig()
		loggerConfig.OutputPaths = []string{config.ProjectRootPath + "/logs/logs.json"}
	}

	logger, err := loggerConfig.Build()
	if err != nil {
		panic(err)
	}

	return logger
}
