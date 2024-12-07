package domain

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"math/rand/v2"

	"github.com/central-university-dev/backend_academy_2024_project_4-go-Dabzelos/pkg"
)

type AffineTransformation struct {
	A, B, C, D, E, F     float64
	TransformationColour color.RGBA
}

type ImageMatrix struct {
	Resolution               *Resolution
	StartingPoints           int
	Pixels                   [][]Pixel
	LinearTransformations    []AffineTransformation
	NonLinearTransformations []func(x, y float64) (newX, newY float64)
}

type Pixel struct {
	X, Y    int
	HitRate int
	Colour  color.RGBA
	normal  float64
}

type Resolution struct {
	Width  int
	Height int
}

const amountOfAffine = 7

func NewImageMatrix(width, height, startingPoints int) *ImageMatrix {
	resolution := Resolution{
		Width:  width,
		Height: height,
	}

	NonlinearTransformations := make([]func(x, y float64) (newX, newY float64), 0)

	matrix := make([][]Pixel, resolution.Height)
	for y := 0; y < resolution.Height; y++ {
		matrix[y] = make([]Pixel, resolution.Width)
		for x := 0; x < resolution.Width; x++ {
			matrix[y][x] = Pixel{X: x, Y: y, Colour: color.RGBA{R: 0, G: 0, B: 0, A: 255}}
		}
	}

	Affine := make([]AffineTransformation, amountOfAffine)

	return &ImageMatrix{Pixels: matrix, Resolution: &resolution, LinearTransformations: Affine,
		NonLinearTransformations: NonlinearTransformations, StartingPoints: startingPoints}
}

func (im *ImageMatrix) GetNonLinearTransform(x, y float64) (newX, newY float64) {
	k, _ := pkg.GenerateRandInt(len(im.NonLinearTransformations))
	return im.NonLinearTransformations[k](x, y)
}

func (im *ImageMatrix) GenerateAffineTransformations() {
	for i := 0; i < amountOfAffine; i++ {
		im.LinearTransformations[i] = generateCoefficients()
	}
}

func (im *ImageMatrix) GetAffineTransform() AffineTransformation {
	x, _ := pkg.GenerateRandInt(amountOfAffine)
	return im.LinearTransformations[x]
}

func generateRandomColor() color.RGBA {
	rCoef, _ := pkg.GenerateRandInt(256)
	gCoef, _ := pkg.GenerateRandInt(256)
	bCoef, _ := pkg.GenerateRandInt(256)
	r := uint8(rCoef)
	g := uint8(gCoef)
	b := uint8(bCoef)

	return color.RGBA{R: r, G: g, B: b, A: 255}
}

func generateCoefficients() AffineTransformation {
	for {
		a := rand.Float64()*2 - 1
		b := rand.Float64()*2 - 1
		d := rand.Float64()*2 - 1
		e := rand.Float64()*2 - 1

		if math.Pow(a, 2)+math.Pow(d, 2) < 1 &&
			math.Pow(b, 2)+math.Pow(e, 2) < 1 &&
			math.Pow(a, 2)+math.Pow(b, 2)+math.Pow(d, 2)+math.Pow(e, 2) < 1+math.Pow(a*e-b*d, 2) {
			fmt.Println("we are here")
			fmt.Println(a, b, d, e)
			c := rand.Float64()*2 - 1
			f := rand.Float64()*2 - 1
			colour := generateRandomColor()

			return AffineTransformation{
				A:                    a,
				B:                    b,
				C:                    c,
				D:                    d,
				E:                    e,
				F:                    f,
				TransformationColour: colour,
			}
		}
	}
}

func AverageColor(c1, c2 color.RGBA) color.RGBA {
	var r, g, b uint8
	red, green, blue, _ := c2.RGBA()
	redNew, greenNew, blueNew, _ := c2.RGBA()

	r, g, b = uint8((red>>8+redNew>>8)/2), uint8((green>>8+greenNew>>8)/2), uint8((blue>>8+blueNew>>8)/2)

	return color.RGBA{
		R: r,
		G: g,
		B: b,
		A: 255,
	}
}

func (im *ImageMatrix) Correction(gamma float64) {
	var maxNormalizedHitRate float64

	for row := range im.Pixels {
		for col := range im.Pixels[row] {
			if im.Pixels[row][col].HitRate != 0 {
				im.Pixels[row][col].normal = math.Log10(float64(im.Pixels[row][col].HitRate))
				if im.Pixels[row][col].normal > maxNormalizedHitRate {
					maxNormalizedHitRate = im.Pixels[row][col].normal
				}
			}
		}
	}

	for row := range im.Pixels {
		for col := range im.Pixels[row] {
			adjusted := math.Pow(im.Pixels[row][col].normal, 1.0/gamma)

			im.Pixels[row][col].Colour.R = uint8(float64(im.Pixels[row][col].Colour.R) * adjusted)
			im.Pixels[row][col].Colour.G = uint8(float64(im.Pixels[row][col].Colour.G) * adjusted)
			im.Pixels[row][col].Colour.B = uint8(float64(im.Pixels[row][col].Colour.B) * adjusted)

		}
	}
}

func (im *ImageMatrix) ReflectHorizontally() {
	for y := 0; y < len(im.Pixels); y++ {
		for x := 0; x < len(im.Pixels[y])/2; x++ {
			mirrorX := len(im.Pixels[y]) - 1 - x
			im.Pixels[y][mirrorX] = im.Pixels[y][x]
		}
	}
}

func (im *ImageMatrix) ReflectVertically() {
	for y := 0; y < len(im.Pixels)/2; y++ {
		mirrorY := len(im.Pixels) - 1 - y
		for x := 0; x < len(im.Pixels[y]); x++ {
			im.Pixels[mirrorY][x] = im.Pixels[y][x]
		}
	}
}

func (im *ImageMatrix) ToImage() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, im.Resolution.Width, im.Resolution.Height))

	for y, row := range im.Pixels {
		for x := range row {
			img.Set(x, y, row[x].Colour)
		}
	}

	return img
}
