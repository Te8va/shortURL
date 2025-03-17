package main

import (
	"go.uber.org/zap"

	"github.com/Te8va/shortURL/internal/app/app"
)

func main() {
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
