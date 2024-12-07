package application

import (
	"encoding/json"
	"fmt"
	"github.com/central-university-dev/backend_academy_2024_project_4-go-Dabzelos/internal/domain"
	"github.com/central-university-dev/backend_academy_2024_project_4-go-Dabzelos/internal/domain/generator"
	"github.com/central-university-dev/backend_academy_2024_project_4-go-Dabzelos/internal/domain/savers"
	"github.com/central-university-dev/backend_academy_2024_project_4-go-Dabzelos/internal/domain/transformations"
	"github.com/central-university-dev/backend_academy_2024_project_4-go-Dabzelos/internal/infrastructure/io"
	"image"
	"log/slog"
	"os"
)

type FractalBuilder interface {
	Render(im *domain.ImageMatrix)
}

type Saver interface {
	Save(img image.Image) error
}

type Config struct {
	Application struct {
		Width              int    `json:"width"`
		Height             int    `json:"height"`
		StartingPoints     int    `json:"startingPoints"`
		SingleThread       bool   `json:"singleThread"`
		Gamma              bool   `json:"gamma"`
		HorizontalSymmetry bool   `json:"horizontalSymmetry"`
		VerticalSymmetry   bool   `json:"verticalSymmetry"`
		Format             string `json:"format"`
	} `json:"Application"`
	LinearTransformations map[string]bool `json:"LinearTransformations"`
}

type Application struct {
	imageMatrix    *domain.ImageMatrix
	symmetry       symmetryFlags
	correction     bool
	outputHandler  *io.Output
	inputHandler   *io.Input
	logger         *slog.Logger
	saver          Saver
	FractalBuilder FractalBuilder
}

type symmetryFlags struct {
	xSymmetry bool
	ySymmetry bool
}

func NewApp(logger *slog.Logger) *Application {
	return &Application{logger: logger}
}

func (a *Application) SetUp(filePath string) error {
	config, err := a.readConfig(filePath)
	if err != nil {
		return fmt.Errorf("ошибка загрузки конфигурации: %w", err)
	}

	a.inputHandler = io.NewReader(os.Stdin)
	a.outputHandler = io.NewWriter(os.Stdout, a.logger)

	a.imageMatrix = domain.NewImageMatrix(config.Application.Width, config.Application.Height, config.Application.StartingPoints)
	a.symmetry = symmetryFlags{
		xSymmetry: config.Application.HorizontalSymmetry,
		ySymmetry: config.Application.VerticalSymmetry,
	}

	a.correction = config.Application.Gamma
	a.saver = a.validateFormat(config.Application.Format)

	a.FractalBuilder = a.validateRenderer(config.Application.SingleThread)

	a.GetSetOfLinearTransformations()

	return nil
}

func (a *Application) readConfig(filePath string) (*Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {

		return nil, err
	}

	return &config, nil
}

func (a *Application) validateFormat(format string) Saver {
	if format == "JPEG" {
		return &savers.JpegSaver{}
	}
	return &savers.PngSaver{}
}

func (a *Application) validateRenderer(singleThread bool) FractalBuilder {
	if singleThread {
		return &generator.SingleThreadGenerator{}
	}
	return &generator.MultiThread{}
}

func (a *Application) GetSetOfLinearTransformations() {
	a.imageMatrix.NonLinearTransformations = append(a.imageMatrix.NonLinearTransformations, transformations.EyeFish)
}

func (a *Application) Start() {
	err := a.SetUp("config.json")
	if err != nil {
		return
	}

	if len(a.imageMatrix.NonLinearTransformations) == 0 {
		a.outputHandler.Write("Unable to generate image without any non linear transformations")
		return
	}

	a.outputHandler.Write("Please wait a bit")
	a.imageMatrix.GenerateAffineTransformations()

	a.FractalBuilder.Render(a.imageMatrix)

	if a.symmetry.xSymmetry {
		a.imageMatrix.ReflectHorizontally()
	}

	if a.symmetry.ySymmetry {
		a.imageMatrix.ReflectVertically()
	}

	if a.correction {
		a.imageMatrix.Correction(2.2)
	}

	img := a.imageMatrix.ToImage()

	err = a.saver.Save(img)
	if err != nil {
		a.outputHandler.Write("Error occurred saving image restart please")
		return
	}

	a.outputHandler.Write("Изображение сохранено как FractalFlame")
}
