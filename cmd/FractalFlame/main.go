package main

import (
	"github.com/central-university-dev/backend_academy_2024_project_4-go-Dabzelos/internal/application"
	"github.com/central-university-dev/backend_academy_2024_project_4-go-Dabzelos/pkg/logger"
)

func main() {
	fileLogger := logger.NewFileLogger("logs.txt")

	defer fileLogger.Close()

	app := application.NewApp(fileLogger.Logger())
	app.Start()
}
