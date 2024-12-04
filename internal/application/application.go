package application

import (
	"log/slog"
	"os"
	"strconv"

	"github.com/central-university-dev/backend_academy_2024_project_4-go-Dabzelos/internal/domain"
	"github.com/central-university-dev/backend_academy_2024_project_4-go-Dabzelos/internal/domain/transformations"
	"github.com/central-university-dev/backend_academy_2024_project_4-go-Dabzelos/internal/infrastructure/io"
)

type Renderer interface {
}

type Saver interface{}

type Application struct {
	imageMatrix   *domain.ImageMatrix
	symmetry      symmetryFlags
	MultiThread   bool
	Correction    bool
	OutputHandler *io.Output
	InputHandler  *io.Input
	logger        *slog.Logger
	format        string
	Saver         Saver
	Renderer      Renderer
}

type symmetryFlags struct {
	xSymmetry bool
	ySymmetry bool
}

func NewApp(logger *slog.Logger) *Application {
	return &Application{logger: logger}
}

func (a *Application) SetUp() {
	a.InputHandler = io.NewReader(os.Stdin)
	a.OutputHandler = io.NewWriter(os.Stdout, a.logger)

	xSymmetry := a.validateFlags("Do you need horizontal symmetry?")
	ySymmetry := a.validateFlags("Do you need vertical symmetry?")

	a.symmetry = symmetryFlags{xSymmetry, ySymmetry}

	thread := a.validateFlags("Do you want multi thread?")

	a.MultiThread = thread

	gammaCor := a.validateFlags("Do you need gamma correction?")

	a.Correction = gammaCor

	xRes, yRes := a.validateResolution()
	imageMatrix := domain.NewImageMatrix(xRes, yRes)

	a.imageMatrix = imageMatrix
}

func (a *Application) validateResolution() (width, height int) {
	a.OutputHandler.Write("Please enter width")

	for {
		userInput, err := a.InputHandler.Read()

		width, err = strconv.Atoi(userInput)
		if err != nil {
			a.OutputHandler.Write("Please enter positive integer")
			continue
		}

		if width <= 0 {
			a.OutputHandler.Write("Please enter positive integer")
			continue
		}

		break
	}

	a.OutputHandler.Write("Please enter height")

	for {
		userInput, err := a.InputHandler.Read()
		height, err = strconv.Atoi(userInput)

		if err != nil {
			a.OutputHandler.Write("Please enter positive integer")
			continue
		}

		if height <= 0 {
			a.OutputHandler.Write("Please enter positive integer")

			continue
		}

		break
	}

	return width, height
}

func (a *Application) validateFlags(message string) (flag bool) {
	a.OutputHandler.Write(message)
	a.OutputHandler.Write("Please enter yes/y or any other input for no")

	userInput, _ := a.InputHandler.Read()
	if userInput == "yes" || userInput == "y" {
		return true
	}

	return false
}

func (a *Application) validateFormat() {}

func (a *Application) validateSetOfLinearTransformations() {
	a.OutputHandler.Write("We need to choose Non Linear transformations")

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

	if a.validateFlags("Do you need  transformation?") {
		a.imageMatrix.NonLinearTransformations = append(a.imageMatrix.NonLinearTransformations, transformations.Spherical)
	}
}

func (a *Application) Start() {
	a.SetUp()
	a.OutputHandler.Write("Please wait a bit")
	a.imageMatrix.GenerateAffineTransformations()

	if a.MultiThread {
		a.imageMatrix.Render()
	} else {

	}

	if a.symmetry.xSymmetry {
		a.imageMatrix.HorizontalSymmetry()
	}

	if a.symmetry.ySymmetry {
		a.imageMatrix.VerticalSymmetry()
	}

	if a.Correction {
		a.imageMatrix.Correction(2.2)
	}

	_ = a.imageMatrix.ToImage()
}
