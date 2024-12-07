package generator

import (
	"math"
	"math/rand/v2"

	"github.com/central-university-dev/backend_academy_2024_project_4-go-Dabzelos/internal/domain"
)

type SingleThreadGenerator struct{}

func (s *SingleThreadGenerator) Render(im *domain.ImageMatrix) {
	var xMin, yMin, xMax, yMax float64

	if im.Resolution.Width > im.Resolution.Height {
		k := float64(im.Resolution.Width) / float64(im.Resolution.Height)
		xMin, yMin, xMax, yMax = -k, -1, k, 1
	} else {
		k := float64(im.Resolution.Height) / float64(im.Resolution.Width)
		xMin, yMin, xMax, yMax = -1, -k, 1, k
	}

	for i := 0; i < im.StartingPoints; i++ {
		newX := rand.Float64()*(xMax-xMin) + xMin
		newY := rand.Float64()*(xMax-xMin) + xMin

		for step := -20; step < 100_000; step++ {
			linearCoeffs := im.GetAffineTransform()
			x := linearCoeffs.A*newX + linearCoeffs.B*newY + linearCoeffs.C
			y := linearCoeffs.D*newY + linearCoeffs.E*newX - linearCoeffs.F

			if step >= 0 {
				pixelX := im.Resolution.Width - int(math.Trunc(((xMax-x)/(xMax-xMin))*float64(im.Resolution.Width)))
				pixelY := im.Resolution.Height - int(math.Trunc(((yMax-y)/(yMax-yMin))*float64(im.Resolution.Height)))

				if pixelX >= 0 && pixelY >= 0 && pixelY < im.Resolution.Height && pixelX < im.Resolution.Width {
					if im.Pixels[pixelY][pixelX].HitRate == 0 {
						im.Pixels[pixelY][pixelX].Colour = linearCoeffs.TransformationColour
						im.Pixels[pixelY][pixelX].HitRate++

						continue
					}

					im.Pixels[pixelY][pixelX].Colour = domain.AverageColor(im.Pixels[pixelY][pixelX].Colour, linearCoeffs.TransformationColour)
					im.Pixels[pixelY][pixelX].HitRate++
				}
			}
			newX, newY = im.GetNonLinearTransform(x, y)
		}
	}
}
