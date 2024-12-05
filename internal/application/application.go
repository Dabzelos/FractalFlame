package application

import (
	"image"
	"log/slog"
	"os"
	"strconv"

	"github.com/central-university-dev/backend_academy_2024_project_4-go-Dabzelos/internal/domain"
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

func (a *Application) SetUp() {
	a.inputHandler = io.NewReader(os.Stdin)
	a.outputHandler = io.NewWriter(os.Stdout, a.logger)

	xSymmetry := a.validateFlags("Do you need horizontal symmetry?")
	ySymmetry := a.validateFlags("Do you need vertical symmetry?")

	a.symmetry = symmetryFlags{xSymmetry, ySymmetry}

	a.validateRenderer()

	gammaCor := a.validateFlags("Do you need gamma correction?")

	a.correction = gammaCor

	xRes, yRes := a.validateResolution()
	samples := a.validateNumberOfStartingPoints()
	imageMatrix := domain.NewImageMatrix(xRes, yRes, samples)

	a.imageMatrix = imageMatrix

	a.GetSetOfLinearTransformations()
	a.validateFormat()
}

func (a *Application) validateRenderer() {
	a.outputHandler.Write("Do you want multi thread generation?")
	a.outputHandler.Write("Please enter yes/y or any other input for no")

	userInput, _ := a.inputHandler.Read()
	if userInput == "yes" || userInput == "y" {
		a.FractalBuilder = &generator.MultiThread{}

		return
	}

	a.FractalBuilder = &generator.SingleThreadGenerator{}
}

func (a *Application) validateNumberOfStartingPoints() int {
	a.outputHandler.Write("Please enter number of starting points: ")

	for {
		userInput, err := a.inputHandler.Read()
		if err != nil {
			a.outputHandler.Write("something went wrong please try again")

			continue
		}

		samples, err := strconv.Atoi(userInput)
		if err != nil {
			a.outputHandler.Write("Please enter positive integer")
			continue
		}

		if samples <= 0 {
			a.outputHandler.Write("Please enter positive integer")
			continue
		}

		return samples
	}
}

func (a *Application) validateResolution() (width, height int) {
	a.outputHandler.Write("Please enter width")

	for {
		userInput, err := a.inputHandler.Read()
		if err != nil {
			a.outputHandler.Write("something went wrong please try again")

			continue
		}

		width, err = strconv.Atoi(userInput)
		if err != nil {
			a.outputHandler.Write("Please enter positive integer")
			continue
		}

		if width <= 0 {
			a.outputHandler.Write("Please enter positive integer")
			continue
		}

		break
	}

	a.outputHandler.Write("Please enter height")

	for {
		userInput, err := a.inputHandler.Read()
		if err != nil {
			a.outputHandler.Write("something went wrong please try again")

			continue
		}

		height, err = strconv.Atoi(userInput)
		if err != nil {
			a.outputHandler.Write("Please enter positive integer")
			continue
		}

		if height <= 0 {
			a.outputHandler.Write("Please enter positive integer")

			continue
		}

		break
	}

	return width, height
}

func (a *Application) validateFlags(message string) (flag bool) {
	a.outputHandler.Write(message)
	a.outputHandler.Write("Please enter yes/y or any other input for no")

	userInput, _ := a.inputHandler.Read()
	if userInput == "yes" || userInput == "y" {
		return true
	}

	return false
}

func (a *Application) validateFormat() {
	a.outputHandler.Write("Please choose format you want to save your picture" +
		"1. PNG\n2. JPEG\n for any non valid input PNG will be set by default\n Please enter PNG or JPEG")

	for {
		userInput, err := a.inputHandler.Read()
		if err != nil {
			a.outputHandler.Write("Something went wrong, please try again")

			continue
		}

		if userInput == "JPEG" {
			a.saver = &savers.JpegSaver{}

			return
		}

		a.saver = &savers.PngSaver{}

		return
	}
}

func (a *Application) GetSetOfLinearTransformations() {
	a.outputHandler.Write("We need to choose Non Linear transformations")

	if a.validateFlags("Do you need sinusoidal transformation?") {
		a.imageMatrix.NonLinearTransformations = append(a.imageMatrix.NonLinearTransformations, transformations.Sinusoidal)
	}

	if a.validateFlags("Do you need disc transformation?") {
		a.imageMatrix.NonLinearTransformations = append(a.imageMatrix.NonLinearTransformations, transformations.Disc)
	}

	if a.validateFlags("Do you need handkerchief transformation?") {
		a.imageMatrix.NonLinearTransformations = append(a.imageMatrix.NonLinearTransformations, transformations.Handkerchief)
	}

	if a.validateFlags("Do you need polar transformation?") {
		a.imageMatrix.NonLinearTransformations = append(a.imageMatrix.NonLinearTransformations, transformations.Polar)
	}

	if a.validateFlags("Do you need horseshoe transformation?") {
		a.imageMatrix.NonLinearTransformations = append(a.imageMatrix.NonLinearTransformations, transformations.Horseshoe)
	}

	if a.validateFlags("Do you need swirl transformation?") {
		a.imageMatrix.NonLinearTransformations = append(a.imageMatrix.NonLinearTransformations, transformations.Swirl)
	}

	if a.validateFlags("Do you need spherical transformation?") {
		a.imageMatrix.NonLinearTransformations = append(a.imageMatrix.NonLinearTransformations, transformations.Spherical)
	}
}

func (a *Application) Start() {
	a.SetUp()

	if len(a.imageMatrix.NonLinearTransformations) == 0 {
		a.outputHandler.Write("Unable to generate image without any non linear transformations")
		return
	}

	a.outputHandler.Write("Please wait a bit")
	a.imageMatrix.GenerateAffineTransformations()

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

	err := a.saver.Save(img)
	if err != nil {
		a.outputHandler.Write("Error occurred saving image restart please")
		return
	}

	a.outputHandler.Write("Изображение сохранено как FractalFlame")
}
