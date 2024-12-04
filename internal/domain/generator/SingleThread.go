package generator

import (
	"github.com/central-university-dev/backend_academy_2024_project_4-go-Dabzelos/internal/domain"
	"math/rand/v2"
)

type SingleThreadGenerator struct{}

func (s *SingleThreadGenerator) Render(im *domain.ImageMatrix) {
	xShapeFactor := float64(im.Resolution.Width) / float64(im.Resolution.Height)
	yShapeFactor := 1.0

	for i := 0; i < im.StartingPoints; i++ {
		newX := (rand.Float64()*2 - 1) * xShapeFactor
		newY := rand.Float64()*2 - 1

		for step := -20; step < 10000; step++ {
			linearCoeffs := im.GetAffineTransform()
			x := linearCoeffs.A*newX + linearCoeffs.B*newY + linearCoeffs.C
			y := linearCoeffs.D*newY + linearCoeffs.E*newX - linearCoeffs.F
			trX, trY := im.GetNonLinearTransform(x, y)

			if step >= 0 { // && trX <= xShapeFactor && trX >= -xShapeFactor && trY <= yShapeFactor && trY >= -yShapeFactor
				pixelX := im.Resolution.Width - int(((xShapeFactor-trX)/(2*xShapeFactor))*float64(im.Resolution.Width))
				pixelY := im.Resolution.Height - int(((yShapeFactor-trY)/(2*yShapeFactor))*float64(im.Resolution.Height))

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
		}
	}
}
