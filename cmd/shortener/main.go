package main

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/Te8va/shortURL/internal/app/app"
)

var (
	buildVersion, buildDate, buildCommit string
)

func main() {
	printGlobalVariable(buildVersion, "version")
	printGlobalVariable(buildDate, "date")
	printGlobalVariable(buildCommit, "commit")

	logger, _ := zap.NewProduction()
	sugar := logger.Sugar()
	defer func() {
		if err := logger.Sync(); err != nil {
			sugar.Errorw("Failed to sync logger", "error", err)
		}
	}()

	appInstance, err := app.NewApp()
	if err != nil {
		sugar.Fatalw("Ошибка инициализации приложения", "error", err)
	}

	if err := appInstance.Run(); err != nil {
		sugar.Fatalw("Ошибка во время работы приложения", "error", err)
	}
}

func printGlobalVariable(variable string, shortDescription string) {
	if variable != "" {
		fmt.Println("Build", shortDescription+":", variable)
	} else {
		fmt.Println("Build", shortDescription+": N/A")
	}
}
