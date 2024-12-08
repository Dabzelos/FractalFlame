package application

import (
	"encoding/json"
	"fmt"

	"image"
	"log/slog"
	"os"

	"github.com/central-university-dev/backend_academy_2024_project_4-go-Dabzelos/internal/domain"
	"github.com/central-university-dev/backend_academy_2024_project_4-go-Dabzelos/internal/domain/errors"
	"github.com/central-university-dev/backend_academy_2024_project_4-go-Dabzelos/internal/domain/generator"
	"github.com/central-university-dev/backend_academy_2024_project_4-go-Dabzelos/internal/domain/savers"
	"github.com/central-university-dev/backend_academy_2024_project_4-go-Dabzelos/internal/domain/transformations"
	"github.com/central-university-dev/backend_academy_2024_project_4-go-Dabzelos/internal/infrastructure/io"
)

type FractalBuilder interface {
	Render(im *domain.ImageMatrix)
}

type Saver interface {
	Save(img image.Image) error
}

type LinearTransformationsConfig struct {
	Spherical    bool `json:"Spherical"`
	Sinusoidal   bool `json:"Sinusoidal"`
	Handkerchief bool `json:"Handkerchief"`
	Swirl        bool `json:"Swirl"`
	Horseshoe    bool `json:"Horseshoe"`
	Polar        bool `json:"Polar"`
	Disc         bool `json:"Disc"`
	Heart        bool `json:"Heart"`
	Linear       bool `json:"Linear"`
	EyeFish      bool `json:"EyeFish"`
}

type Config struct {
	Application struct {
		Width              int    `json:"width"`
		Height             int    `json:"height"`
		StartingPoints     int    `json:"startingPoints"`
		Iterations         int    `json:"iterations"`
		SingleThread       bool   `json:"singleThread"`
		Gamma              bool   `json:"gamma"`
		HorizontalSymmetry bool   `json:"horizontalSymmetry"`
		VerticalSymmetry   bool   `json:"verticalSymmetry"`
		Format             string `json:"format"`
	} `json:"Application"`
	ListOfTransformations LinearTransformationsConfig `json:"LinearTransformations"`
}

type Application struct {
	imageMatrix    *domain.ImageMatrix
	symmetry       symmetryFlags
	correction     bool
	outputHandler  *io.Output
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

	a.outputHandler = io.NewWriter(os.Stdout, a.logger)

	a.imageMatrix = domain.NewImageMatrix(config.Application.Width, config.Application.Height,
		config.Application.StartingPoints, config.Application.Iterations)
	a.symmetry = symmetryFlags{
		xSymmetry: config.Application.HorizontalSymmetry,
		ySymmetry: config.Application.VerticalSymmetry,
	}

	a.correction = config.Application.Gamma
	a.saver = a.validateFormat(config.Application.Format)

	a.FractalBuilder = a.validateRenderer(config.Application.SingleThread)

	a.GetSetOfLinearTransformations(config.ListOfTransformations)

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

	if config.Application.Height == 0 || config.Application.Width == 0 || config.Application.StartingPoints == 0 {
		return nil, errors.ErrZeroSizeMatrix{}
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

	return &generator.MultiThreadGenerator{}
}

func (a *Application) GetSetOfLinearTransformations(trConfig LinearTransformationsConfig) {
	var functions []func(float64, float64) (float64, float64)

	if trConfig.Spherical {
		functions = append(functions, transformations.Spherical)
	}

	if trConfig.Sinusoidal {
		functions = append(functions, transformations.Sinusoidal)
	}

	if trConfig.Handkerchief {
		functions = append(functions, transformations.Handkerchief)
	}

	if trConfig.Swirl {
		functions = append(functions, transformations.Swirl)
	}

	if trConfig.Horseshoe {
		functions = append(functions, transformations.Horseshoe)
	}

	if trConfig.Polar {
		functions = append(functions, transformations.Polar)
	}

	if trConfig.Disc {
		functions = append(functions, transformations.Disc)
	}

	if trConfig.Heart {
		functions = append(functions, transformations.Heart)
	}

	if trConfig.Linear {
		functions = append(functions, transformations.Linear)
	}

	if trConfig.EyeFish {
		functions = append(functions, transformations.EyeFish)
	}

	a.imageMatrix.NonLinearTransformations = functions
}

func (a *Application) Start() {
	if err := a.SetUp("config.json"); err != nil {
		a.logger.Error("Error occurred reading from config file")
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

	img := a.imageMatrix.ConvertToImage()

	if err := a.saver.Save(img); err != nil {
		a.logger.Error("error occurred saving image", "error", err)
		a.outputHandler.Write("Error occurred saving image restart please")

		return
	}

	a.outputHandler.Write("Изображение сохранено как FractalFlame")
}
