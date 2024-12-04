package main

import (
	"fmt"
	"image/png"
	"os"
	"time"

	"github.com/central-university-dev/backend_academy_2024_project_4-go-Dabzelos/internal/domain"
	"github.com/central-university-dev/backend_academy_2024_project_4-go-Dabzelos/internal/domain/transformations"
)

/*
	fileLogger := logger.NewFileLogger("logs.txt")

	defer fileLogger.Close()

	app := application.NewApp(fileLogger.Logger())
	app.Start()*/

func main() {
	start := time.Now()
	ImageMatrix := domain.NewImageMatrix(1920, 1080)

	ImageMatrix.GenerateAffineTransformations()

	ImageMatrix.NonLinearTransformations = append(ImageMatrix.NonLinearTransformations, transformations.Disc, transformations.Sinusoidal)

	ImageMatrix.Render()
	/* ImageMatrix.HorizontalSymmetry()
	ImageMatrix.VerticalSymmetry()*/
	ImageMatrix.Correction(2.2)
	img := ImageMatrix.ToImage()

	file, err := os.Create("spherical_transform.png")
	if err != nil {
		panic(err)
	}
	defer file.Close() // Задаем размеры изображения

	// Генерируем изображение

	err = png.Encode(file, img)
	if err != nil {
		panic(err)
	}

	println("Изображение сохранено как spherical_transform.png")

	elapsed := time.Since(start) // Вычисляем разницу

	minutes := int(elapsed.Minutes())
	seconds := int(elapsed.Seconds()) % 60

	fmt.Printf("Время выполнения: %02d:%02d\n", minutes, seconds)
}
