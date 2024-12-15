package domain

import (
	"image"
	"image/color"
	"math"
	"math/rand/v2"
	"sync"

	"github.com/central-university-dev/backend_academy_2024_project_4-go-Dabzelos/pkg/random"
)

type AffineTransformation struct {
	A, B, C, D, E, F     float64
	TransformationColour color.RGBA
}

type TransformFunc func(x, y float64) (newX, newY float64)

type ImageMatrix struct {
	Resolution               *Resolution
	cords                    CoordinatesRange
	StartingPoints           int
	Iterations               int
	Pixels                   [][]Pixel
	LinearTransformations    []AffineTransformation
	NonLinearTransformations []TransformFunc
}

type Pixel struct {
	X, Y    int
	HitRate int
	Colour  color.RGBA
	normal  float64
	mutex   sync.Mutex
}

type Resolution struct {
	Width  int
	Height int
}

type CoordinatesRange struct {
	xMin, yMin, xMax, yMax float64
}

const amountOfAffine = 7

func NewImageMatrix(width, height, startingPoints, iterations int) *ImageMatrix {
	resolution := Resolution{
		Width:  width,
		Height: height,
	}

	NonlinearTransformations := make([]TransformFunc, 0, 10)

	matrix := make([][]Pixel, resolution.Height)
	for y := 0; y < resolution.Height; y++ {
		matrix[y] = make([]Pixel, resolution.Width)
		for x := 0; x < resolution.Width; x++ {
			matrix[y][x] = Pixel{X: x, Y: y, Colour: color.RGBA{R: 0, G: 0, B: 0, A: 255}}
		}
	}

	var xMin, yMin, xMax, yMax float64

	if width > height {
		k := float64(width) / float64(height)
		xMin, yMin, xMax, yMax = -k, -1, k, 1
	} else {
		k := float64(height) / float64(width)
		xMin, yMin, xMax, yMax = -1, -k, 1, k
	}

	cords := CoordinatesRange{
		xMin: xMin,
		yMin: yMin,
		xMax: xMax,
		yMax: yMax,
	}

	Affine := make([]AffineTransformation, amountOfAffine)

	return &ImageMatrix{Pixels: matrix, Resolution: &resolution, LinearTransformations: Affine,
		NonLinearTransformations: NonlinearTransformations, StartingPoints: startingPoints, Iterations: iterations, cords: cords}
}

// GetNonLinearTransform - возвращает применение к координатам случайной функции нелинейного преобразования.
func (im *ImageMatrix) GetNonLinearTransform(x, y float64) (newX, newY float64) {
	k := rand.IntN(len(im.NonLinearTransformations)) //nolint
	return im.NonLinearTransformations[k](x, y)
}

// GenerateAffineTransformations - функция, которая генерирует все 7(определенно константой) случайных аффинных
// преобразований.
func (im *ImageMatrix) GenerateAffineTransformations() {
	for i := 0; i < amountOfAffine; i++ {
		im.LinearTransformations[i] = im.generateCoefficients()
	}
}

// GetAffineTransform - позволяет получить одно случайное
// из 7(определенно константой) линейных(аффинных) преобразований.
func (im *ImageMatrix) GetAffineTransform() AffineTransformation {
	x := rand.IntN(amountOfAffine) //nolint
	return im.LinearTransformations[x]
}

// generateCoefficients -позволяет сгенерировать коэффициенты и цвет для линейного преобразования.
func (im *ImageMatrix) generateCoefficients() AffineTransformation {
	for {
		a := random.GenerateRandFloat64()
		b := random.GenerateRandFloat64()
		d := random.GenerateRandFloat64()
		e := random.GenerateRandFloat64()

		if math.Pow(a, 2)+math.Pow(d, 2) < 1 &&
			math.Pow(b, 2)+math.Pow(e, 2) < 1 &&
			math.Pow(a, 2)+math.Pow(b, 2)+math.Pow(d, 2)+math.Pow(e, 2) < 1+math.Pow(a*e-b*d, 2) {
			c := random.GenerateRandFloat64()
			f := random.GenerateRandFloat64()
			colour := random.GenerateRandomColor()

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

// averageColor - позволяет усреднить два цвета, и вернуть результат усреднения.
func (im *ImageMatrix) averageColor(c1, c2 color.RGBA) color.RGBA {
	var r, g, b uint8

	red, green, blue, _ := c1.RGBA()
	redNew, greenNew, blueNew, _ := c2.RGBA()

	r, g, b = byte((red>>8+redNew>>8)>>1), byte((green>>8+greenNew>>8)>>1), byte((blue>>8+blueNew>>8)>>1)

	return color.RGBA{
		R: r,
		G: g,
		B: b,
		A: 255,
	}
}

// Correction - реализация алгоритма гамма коррекции.
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

// ReflectHorizontally - реализация симметрии по иксу.
func (im *ImageMatrix) ReflectHorizontally() {
	for y := 0; y < len(im.Pixels); y++ {
		for x := 0; x < len(im.Pixels[y])/2; x++ {
			mirrorX := len(im.Pixels[y]) - 1 - x
			im.Pixels[y][mirrorX] = Pixel{
				Colour:  im.Pixels[y][x].Colour,
				X:       x,
				Y:       y,
				HitRate: im.Pixels[y][x].HitRate,
			}
		}
	}
}

// ReflectVertically - реализация симметрии по игреку.
func (im *ImageMatrix) ReflectVertically() {
	for y := 0; y < len(im.Pixels)/2; y++ {
		mirrorY := len(im.Pixels) - 1 - y

		copy(im.Pixels[mirrorY], im.Pixels[y])
	}
}

// ConvertToImage - преобразовывает структуру ImageMatrix в картинку.
func (im *ImageMatrix) ConvertToImage() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, im.Resolution.Width, im.Resolution.Height))

	for y, row := range im.Pixels {
		for x := range row {
			img.Set(x, y, row[x].Colour)
		}
	}

	return img
}

// UpdatePixel - отвечает за обработку одного пикселя в рамках работы алгоритма.
func (im *ImageMatrix) UpdatePixel(pixelY, pixelX int, linearCoeffs AffineTransformation) {
	im.Pixels[pixelY][pixelX].mutex.Lock()
	defer im.Pixels[pixelY][pixelX].mutex.Unlock()

	if im.Pixels[pixelY][pixelX].HitRate == 0 {
		im.Pixels[pixelY][pixelX].Colour = linearCoeffs.TransformationColour
		im.Pixels[pixelY][pixelX].HitRate++

		return
	}

	im.Pixels[pixelY][pixelX].Colour = im.averageColor(im.Pixels[pixelY][pixelX].Colour, linearCoeffs.TransformationColour)
	im.Pixels[pixelY][pixelX].HitRate++
}

// GenerateStartingCoordinates - позволяет получить координаты стартовых точек для работы алгоритма.
func (im *ImageMatrix) GenerateStartingCoordinates() (newX, newY float64) {
	newX = random.GenerateRandFloat64() //nolint
	newY = random.GenerateRandFloat64() //nolint

	newX = newX*(im.cords.xMax-im.cords.xMin) + im.cords.xMin
	newY = newY*(im.cords.yMax-im.cords.yMin) + im.cords.yMin

	return newX, newY
}

// ProcessStartingPoint - функция реализующая логику обработки каждой стартовой точки, вынесено в отдельную во избежание
// дублирования кода.
func (im *ImageMatrix) ProcessStartingPoint() {
	newX, newY := im.GenerateStartingCoordinates()

	for step := -20; step < im.Iterations; step++ {
		linearCoeffs := im.GetAffineTransform() // Получаем линейные коэффициенты трансформации
		x := linearCoeffs.A*newX + linearCoeffs.B*newY + linearCoeffs.C
		y := linearCoeffs.D*newY + linearCoeffs.E*newX - linearCoeffs.F

		if step >= 0 {
			pixelX := im.Resolution.Width - int(math.Trunc(((im.cords.xMax-x)/(im.cords.xMax-im.cords.xMin))*
				float64(im.Resolution.Width)))
			pixelY := im.Resolution.Height - int(math.Trunc(((im.cords.yMax-y)/(im.cords.yMax-im.cords.yMin))*
				float64(im.Resolution.Height)))

			if pixelX >= 0 && pixelY >= 0 && pixelY < im.Resolution.Height && pixelX < im.Resolution.Width {
				im.UpdatePixel(pixelY, pixelX, linearCoeffs)
			}
		}

		newX, newY = im.GetNonLinearTransform(x, y)
	}
}
