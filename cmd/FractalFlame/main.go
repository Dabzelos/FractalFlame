package main

import (
	"flag"
	"os"

	"github.com/central-university-dev/backend_academy_2024_project_4-go-Dabzelos/internal/application"
	"github.com/central-university-dev/backend_academy_2024_project_4-go-Dabzelos/internal/infrastructure/io"
	"github.com/central-university-dev/backend_academy_2024_project_4-go-Dabzelos/pkg/logger"
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
