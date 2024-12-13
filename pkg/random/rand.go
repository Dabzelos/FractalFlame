package random

import (
	"image/color"
	"math/rand/v2"
)

// GenerateRandFloat64 позволяет получить значение в формате типа float64, из диапазона [-1;1].
func GenerateRandFloat64() float64 {
	n := rand.Float64() //nolint

	return n*2 - 1
}

// GenerateRandomColor - функция, которая генерирует случайный цвет в цветовой модели RGBA.
func GenerateRandomColor() color.RGBA {
	return color.RGBA{
		R: byte(rand.IntN(255)), //nolint
		G: byte(rand.IntN(255)), //nolint
		B: byte(rand.IntN(255)), //nolint
		A: 255}
}
