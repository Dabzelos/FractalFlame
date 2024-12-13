package generator_test

import (
	"fmt"

	"testing"

	"github.com/central-university-dev/backend_academy_2024_project_4-go-Dabzelos/internal/domain"
	"github.com/central-university-dev/backend_academy_2024_project_4-go-Dabzelos/internal/domain/generator"
	"github.com/central-university-dev/backend_academy_2024_project_4-go-Dabzelos/internal/domain/transformations"
)

func BenchmarkSingleThreadGenerator_Render(b *testing.B) {
	tc := []struct {
		width          int
		height         int
		StartingPoints int
		Iterations     int
	}{
		{width: 1980, height: 1080, StartingPoints: 100, Iterations: 100000},
		{width: 1980, height: 1080, StartingPoints: 150, Iterations: 100000},
		{width: 1980, height: 1080, StartingPoints: 200, Iterations: 100000},
		{width: 2560, height: 1440, StartingPoints: 100, Iterations: 100000},
		{width: 2560, height: 1440, StartingPoints: 150, Iterations: 100000},
		{width: 2560, height: 1440, StartingPoints: 200, Iterations: 100000},
	}

	for _, tt := range tc {
		b.Run(fmt.Sprintf("width: %d, height: %d, Starting points %d", tt.width, tt.height,
			tt.StartingPoints), func(b *testing.B) {
			img := domain.NewImageMatrix(tt.width, tt.height, tt.StartingPoints, tt.Iterations)
			gn := &generator.SingleThreadGenerator{}

			img.GenerateAffineTransformations()

			img.NonLinearTransformations = append(img.NonLinearTransformations, transformations.Disc,
				transformations.Linear, transformations.Polar)

			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				gn.Render(img)
			}
		})
	}
}

func BenchmarkMultiThreadGenerator_Render(b *testing.B) {
	tc := []struct {
		width          int
		height         int
		StartingPoints int
		Iterations     int
		NumWorkers     int
	}{
		{width: 1980, height: 1080, StartingPoints: 100, Iterations: 100000, NumWorkers: 8},
		{width: 1980, height: 1080, StartingPoints: 150, Iterations: 100000, NumWorkers: 8},
		{width: 1980, height: 1080, StartingPoints: 200, Iterations: 100000, NumWorkers: 8},
		{width: 2560, height: 1440, StartingPoints: 100, Iterations: 100000, NumWorkers: 8},
		{width: 2560, height: 1440, StartingPoints: 150, Iterations: 100000, NumWorkers: 8},
		{width: 2560, height: 1440, StartingPoints: 200, Iterations: 100000, NumWorkers: 8},
	}

	for _, tt := range tc {
		b.Run(fmt.Sprintf("width: %d, height: %d, Starting points %d, numWorkers %d", tt.width, tt.height,
			tt.StartingPoints, tt.NumWorkers), func(b *testing.B) {
			img := domain.NewImageMatrix(tt.width, tt.height, tt.StartingPoints, tt.Iterations)
			gn := &generator.MultiThreadGenerator{NumWorkers: tt.NumWorkers}

			img.GenerateAffineTransformations()
			img.NonLinearTransformations = append(img.NonLinearTransformations, transformations.Disc,
				transformations.Linear, transformations.Polar)

			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				gn.Render(img)
			}
		})
	}
}

func BenchmarkMultiThreadGenerator_Render16(b *testing.B) {
	tc := []struct {
		width          int
		height         int
		StartingPoints int
		Iterations     int
		NumWorkers     int
	}{
		{width: 1980, height: 1080, StartingPoints: 100, Iterations: 100000, NumWorkers: 16},
		{width: 1980, height: 1080, StartingPoints: 150, Iterations: 100000, NumWorkers: 16},
		{width: 1980, height: 1080, StartingPoints: 200, Iterations: 100000, NumWorkers: 16},
		{width: 2560, height: 1440, StartingPoints: 100, Iterations: 100000, NumWorkers: 16},
		{width: 2560, height: 1440, StartingPoints: 150, Iterations: 100000, NumWorkers: 16},
		{width: 2560, height: 1440, StartingPoints: 200, Iterations: 100000, NumWorkers: 16},
	}

	for _, tt := range tc {
		b.Run(fmt.Sprintf("width: %d, height: %d, Starting points %d, numWorkers %d", tt.width, tt.height,
			tt.StartingPoints, tt.NumWorkers), func(b *testing.B) {
			img := domain.NewImageMatrix(tt.width, tt.height, tt.StartingPoints, tt.Iterations)
			gn := &generator.MultiThreadGenerator{NumWorkers: tt.NumWorkers}

			img.GenerateAffineTransformations()
			img.NonLinearTransformations = append(img.NonLinearTransformations, transformations.Disc,
				transformations.Linear, transformations.Polar)

			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				gn.Render(img)
			}
		})
	}
}

func BenchmarkMultiThreadGenerator_Render24(b *testing.B) {
	tc := []struct {
		width          int
		height         int
		StartingPoints int
		Iterations     int
		NumWorkers     int
	}{
		{width: 1980, height: 1080, StartingPoints: 100, Iterations: 100000, NumWorkers: 24},
		{width: 1980, height: 1080, StartingPoints: 150, Iterations: 100000, NumWorkers: 24},
		{width: 1980, height: 1080, StartingPoints: 200, Iterations: 100000, NumWorkers: 24},
		{width: 2560, height: 1440, StartingPoints: 100, Iterations: 100000, NumWorkers: 24},
		{width: 2560, height: 1440, StartingPoints: 150, Iterations: 100000, NumWorkers: 24},
		{width: 2560, height: 1440, StartingPoints: 200, Iterations: 100000, NumWorkers: 24},
	}

	for _, tt := range tc {
		b.Run(fmt.Sprintf("width: %d, height: %d, Starting points %d, numWorkers %d", tt.width, tt.height,
			tt.StartingPoints, tt.NumWorkers), func(b *testing.B) {
			img := domain.NewImageMatrix(tt.width, tt.height, tt.StartingPoints, tt.Iterations)
			gn := &generator.MultiThreadGenerator{NumWorkers: tt.NumWorkers}

			img.GenerateAffineTransformations()
			img.NonLinearTransformations = append(img.NonLinearTransformations, transformations.Disc,
				transformations.Linear, transformations.Polar)

			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				gn.Render(img)
			}
		})
	}
}
