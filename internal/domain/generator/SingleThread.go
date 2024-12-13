package generator

import (
	"github.com/central-university-dev/backend_academy_2024_project_4-go-Dabzelos/internal/domain"
)

type SingleThreadGenerator struct{}

// Render функция, которая обеспечивает генерацию фрактального пламени.
func (s *SingleThreadGenerator) Render(im *domain.ImageMatrix) {
	for i := 0; i < im.StartingPoints; i++ {
		im.ProcessStartingPoint()
	}
}
