package application

import (
	"image"
	"log/slog"
	"os"

	"FractalFlame/configuration"
	"FractalFlame/internal/domain"
	"FractalFlame/internal/domain/errors"
	"FractalFlame/internal/domain/generator"
	"FractalFlame/internal/domain/savers"
	"FractalFlame/internal/domain/transformations"
	"FractalFlame/internal/infrastructure/io"
)

type fractalBuilder interface {
	Render(im *domain.ImageMatrix)
}

type saver interface {
	Save(img image.Image) error
}

type outputHandler interface {
	Write(messages ...interface{})
}

type Application struct {
	imageMatrix     *domain.ImageMatrix
	symmetry        symmetryFlags
	correction      bool
	correctionCoeff float64
	outputHandler   outputHandler
	logger          *slog.Logger
	saver           saver
	fractalBuilder  fractalBuilder
}

type symmetryFlags struct {
	xSymmetry bool
	ySymmetry bool
}

func NewApp(logger *slog.Logger, handler outputHandler) *Application {
	return &Application{logger: logger, outputHandler: handler}
}

func (a *Application) setUp(filePath *string) error {
	config, err := configuration.Read(*filePath)
	if err != nil {
		return errors.ErrReadingConfig{Err: err}
	}

	a.outputHandler = io.NewWriter(os.Stdout, a.logger)

	a.imageMatrix = domain.NewImageMatrix(config.Application.Width, config.Application.Height,
		config.Application.StartingPoints, config.Application.Iterations)
	a.symmetry = symmetryFlags{
		xSymmetry: config.Application.HorizontalSymmetry,
		ySymmetry: config.Application.VerticalSymmetry,
	}

	a.correction = config.Application.Gamma
	a.correctionCoeff = config.Application.GammaCoeff
	a.setSaver(config.Application.Format)
	a.setRenderer(config.Application.SingleThread, config.Application.NumWorkers)
	a.validateSetOfLinearTransformations(config.ListOfTransformations)

	return nil
}

func (a *Application) setSaver(format string) {
	if format == "JPEG" {
		a.saver = &savers.JpegSaver{}

		return
	}

	a.saver = &savers.PngSaver{}
}

func (a *Application) setRenderer(singleThread bool, workers int) {
	if singleThread {
		a.fractalBuilder = &generator.SingleThreadGenerator{}

		return
	}

	a.fractalBuilder = &generator.MultiThreadGenerator{NumWorkers: workers}
}

func (a *Application) validateSetOfLinearTransformations(trConfig configuration.LinearTransformationsConfig) {
	var functions []domain.TransformFunc

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

func (a *Application) Start(source *string) error {
	if err := a.setUp(source); err != nil {
		return errors.ErrReadingConfig{Err: err}
	}

	if len(a.imageMatrix.NonLinearTransformations) == 0 {
		a.outputHandler.Write("Unable to generate image without any non linear transformations")
		return errors.ErrZeroSizeMatrix{}
	}

	a.imageMatrix.GenerateAffineTransformations()

	a.fractalBuilder.Render(a.imageMatrix)

	if a.symmetry.xSymmetry {
		a.imageMatrix.ReflectHorizontally()
	}

	if a.symmetry.ySymmetry {
		a.imageMatrix.ReflectVertically()
	}

	if a.correction {
		a.imageMatrix.Correction(a.correctionCoeff)
	}

	img := a.imageMatrix.ConvertToImage()

	if err := a.saver.Save(img); err != nil {
		a.outputHandler.Write("Error occurred saving image restart please")

		return errors.ErrSavingImage{Err: err}
	}

	a.outputHandler.Write("Изображение сохранено как FractalFlame")

	return nil
}
