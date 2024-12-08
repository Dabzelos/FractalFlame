package generator

import (
	"runtime"
	"sync"

	"github.com/central-university-dev/backend_academy_2024_project_4-go-Dabzelos/internal/domain"
)

type MultiThreadGenerator struct {
}

// Render функция, которая обеспечивает многопоточную генерацию фрактального пламени.
func (m *MultiThreadGenerator) Render(im *domain.ImageMatrix) {
	var xMin, yMin, xMax, yMax float64

	if im.Resolution.Width > im.Resolution.Height {
		k := float64(im.Resolution.Width) / float64(im.Resolution.Height)
		xMin, yMin, xMax, yMax = -k, -1, k, 1
	} else {
		k := float64(im.Resolution.Height) / float64(im.Resolution.Width)
		xMin, yMin, xMax, yMax = -1, -k, 1, k
	}

	jobs := make(chan int, im.StartingPoints)

	var wg sync.WaitGroup

	numWorkers := runtime.NumCPU()

	for w := 0; w < numWorkers; w++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for range jobs {
				im.ProcessStartingPoint(xMin, yMin, xMax, yMax)
			}
		}()
	}

	for i := 0; i < im.StartingPoints; i++ {
		jobs <- i
	}

	close(jobs)

	wg.Wait()
}
