package generator

import (
	"runtime"
	"sync"

	"github.com/central-university-dev/backend_academy_2024_project_4-go-Dabzelos/internal/domain"
)

type MultiThreadGenerator struct {
	NumWorkers int
}

// Render функция, которая обеспечивает многопоточную генерацию фрактального пламени.
func (m *MultiThreadGenerator) Render(im *domain.ImageMatrix) {
	var wg sync.WaitGroup

	jobs := make(chan int, im.StartingPoints)

	if m.NumWorkers == 0 {
		m.NumWorkers = runtime.NumCPU()
	}

	for w := 0; w < m.NumWorkers; w++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for range jobs {
				im.ProcessStartingPoint()
			}
		}()
	}

	for i := 0; i < im.StartingPoints; i++ {
		jobs <- i
	}

	close(jobs)

	wg.Wait()
}
