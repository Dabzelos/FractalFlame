package generator

import (
	"github.com/central-university-dev/backend_academy_2024_project_4-go-Dabzelos/internal/domain"
)

type SingleThreadGenerator struct{}

// Render функция, которая обеспечивает генерацию фрактального пламени.
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
		im.ProcessStartingPoint(xMin, yMin, xMax, yMax)
	}
}
