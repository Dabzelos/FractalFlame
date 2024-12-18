package main

import (
	"flag"
	"os"

	"FractalFlame/internal/application"
	"FractalFlame/internal/infrastructure/io"
	"FractalFlame/pkg/logger"
)

func main() {
	config := flag.String("config", "", "path")
	flag.Parse()

	fileLogger := logger.NewFileLogger("logs.txt")
	outputHandler := io.NewWriter(os.Stdout, fileLogger.Logger())

	defer fileLogger.Close()

	app := application.NewApp(fileLogger.Logger(), outputHandler)
	if err := app.Start(config); err != nil {
		fileLogger.Logger().Error("Error happened while running the application", "error", err)
	}
}
