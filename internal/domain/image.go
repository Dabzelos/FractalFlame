package domain

import (
	"image"
	"image/color"
	"math"
	"math/rand"

	"github.com/central-university-dev/backend_academy_2024_project_4-go-Dabzelos/pkg"
)

type AffineTransformation struct {
	a, b, c, d, e, f     float64
	TransformationColour color.RGBA
}

type ImageMatrix struct {
	Resolution               *Resolution
	Pixels                   [][]Pixel
	LinearTransformations    []AffineTransformation
	NonLinearTransformations []func(x, y float64) (newX, newY float64)
}

type Pixel struct {
	X, Y    int
	hitRate int
	colour  color.RGBA
	normal  float64
}

type Resolution struct {
	Width  int
	Height int
}

const (
	StartingPoints = 10000
	amountOfAffine = 10
)

func NewImageMatrix(width, height int) *ImageMatrix {
	resolution := Resolution{
		Width:  width,
		Height: height,
	}

	matrix := make([][]Pixel, resolution.Height)
	for y := 0; y < resolution.Height; y++ {
		matrix[y] = make([]Pixel, resolution.Width)
		for x := 0; x < resolution.Width; x++ {
			matrix[y][x] = Pixel{X: x, Y: y, colour: color.RGBA{R: 0, G: 0, B: 0, A: 255}}
		}
	}

	Affine := make([]AffineTransformation, amountOfAffine)

	return &ImageMatrix{Pixels: matrix, Resolution: &resolution, LinearTransformations: Affine}
}

func (im *ImageMatrix) getNonLinearTransform(x, y float64) (newX, newY float64) {
	k, _ := pkg.GenerateRandInt(len(im.NonLinearTransformations))
	return im.NonLinearTransformations[k](x, y)
}

func (im *ImageMatrix) GenerateAffineTransformations() {
	for i := 0; i < amountOfAffine; i++ {
		im.LinearTransformations[i] = generateCoefficients()
	}
}

func (im *ImageMatrix) getAffineTransform() AffineTransformation {
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
		a := rand.Float64()*3 - 1.5
		b := rand.Float64()*3 - 1.5
		d := rand.Float64()*3 - 1.5
		e := rand.Float64()*3 - 1.5

		if math.Pow(a, 2)+math.Pow(d, 2) < 1 &&
			math.Pow(b, 2)+math.Pow(e, 2) < 1 &&
			math.Pow(a, 2)+math.Pow(b, 2)+math.Pow(d, 2)+math.Pow(e, 2) < 1+math.Pow(a*e-b*d, 2) {
			c := rand.Float64()*2 - 1
			f := rand.Float64()*2 - 1
			colour := generateRandomColor()

			return AffineTransformation{
				a:                    a,
				b:                    b,
				c:                    c,
				d:                    d,
				e:                    e,
				f:                    f,
				TransformationColour: colour,
			}
		}
	}
}

func averageColor(c1, c2 color.RGBA) color.RGBA {
	return color.RGBA{
		R: ((c1.R) + (c2.R)) / 2,
		G: ((c1.G) + (c2.G)) / 2,
		B: ((c1.B) + (c2.B)) / 2,
		A: ((c1.A) + (c2.A)) / 2,
	}
}

func (im *ImageMatrix) Render() {
	xShapeFactor := float64(im.Resolution.Width) / float64(im.Resolution.Height)
	yShapeFactor := 1.0

	for i := 0; i < StartingPoints; i++ {
		newX := (rand.Float64()*2 - 1) * xShapeFactor
		newY := rand.Float64()*2 - 1

		for step := -20; step < 1_000; step++ {
			linearCoeffs := im.getAffineTransform()
			x := linearCoeffs.a*newX + linearCoeffs.b*newY + linearCoeffs.c
			y := linearCoeffs.d*newY + linearCoeffs.e*newX - linearCoeffs.f
			trX, trY := im.getNonLinearTransform(x, y)

			if step >= 0 { // && trX <= xShapeFactor && trX >= -xShapeFactor && trY <= yShapeFactor && trY >= -yShapeFactor
				pixelX := im.Resolution.Width - int(((xShapeFactor-trX)/(2*xShapeFactor))*float64(im.Resolution.Width))
				pixelY := im.Resolution.Height - int(((yShapeFactor-trY)/(2*yShapeFactor))*float64(im.Resolution.Height))

				if pixelX >= 0 && pixelY >= 0 && pixelY < im.Resolution.Height && pixelX < im.Resolution.Width {
					if im.Pixels[pixelY][pixelX].hitRate == 0 {
						im.Pixels[pixelY][pixelX].colour = linearCoeffs.TransformationColour
						im.Pixels[pixelY][pixelX].hitRate++

						continue
					}

					im.Pixels[pixelY][pixelX].colour = averageColor(im.Pixels[pixelY][pixelX].colour, linearCoeffs.TransformationColour)
					im.Pixels[pixelY][pixelX].hitRate++
				}
			}
		}
	}
}

func (im *ImageMatrix) Correction(gamma float64) {
	var maxNormalizedHitRate float64

	for row := range im.Pixels {
		for col := range im.Pixels[row] {
			if im.Pixels[row][col].hitRate != 0 {
				im.Pixels[row][col].normal = math.Log10(float64(im.Pixels[row][col].hitRate))
				if im.Pixels[row][col].normal > maxNormalizedHitRate {
					maxNormalizedHitRate = im.Pixels[row][col].normal
				}
			}
		}
	}

	for row := range im.Pixels {
		for col := range im.Pixels[row] {
			adjusted := math.Pow(im.Pixels[row][col].normal, 1.0/gamma)

			// im.Pixels[row][col].colour.A = uint8(float64(im.Pixels[row][col].colour.A) * adjusted)
			im.Pixels[row][col].colour.R = uint8(float64(im.Pixels[row][col].colour.R) * adjusted)
			im.Pixels[row][col].colour.G = uint8(float64(im.Pixels[row][col].colour.G) * adjusted)
			im.Pixels[row][col].colour.B = uint8(float64(im.Pixels[row][col].colour.B) * adjusted)
		}
	}
}

func (im *ImageMatrix) HorizontalSymmetry() {
	for y := 0; y < len(im.Pixels); y++ {
		for x := 0; x < len(im.Pixels[y])/2; x++ {
			// Отражаем пиксели справа налево
			mirrorX := len(im.Pixels[y]) - 1 - x
			im.Pixels[y][mirrorX] = im.Pixels[y][x]
		}
	}
}

func (im *ImageMatrix) VerticalSymmetry() {
	for y := 0; y < len(im.Pixels)/2; y++ {
		mirrorY := len(im.Pixels) - 1 - y
		for x := 0; x < len(im.Pixels[y]); x++ {
			// Отражаем пиксели снизу вверх
			im.Pixels[mirrorY][x] = im.Pixels[y][x]
		}
	}
}

func (im *ImageMatrix) ToImage() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, im.Resolution.Width, im.Resolution.Height))

	for y, row := range im.Pixels {
		for x := range row {
			img.Set(x, y, row[x].colour)
		}
	}

	return img
}
