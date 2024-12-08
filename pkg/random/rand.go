package random

import (
	"crypto/rand"
	"image/color"
	"math/big"
)

// GenerateRandInt - функция позволяющая удобно получать случайное число.
// Вынесено в отдельную функцию для удобства обработки ошибки, а так же тк оба алгоритма многократно пользуются
// генерацией случайного числа.
func GenerateRandInt(maxNum int) (int, error) {
	MaxRand := big.NewInt(int64(maxNum))

	randomIndex, err := rand.Int(rand.Reader, MaxRand)
	if err != nil {
		return int(randomIndex.Int64()), err
	}

	return int(randomIndex.Int64()), nil
}

func GenerateRandFloat64() (float64, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(1000000000))

	if err != nil {
		return 0, err
	}

	return float64(n.Int64())/500000000.0 - 1, nil
}

// GenerateRandomColor - функция, которая генерирует случайный цвет в цветовой модели RGBA.
func GenerateRandomColor() color.RGBA {
	rCoef, _ := GenerateRandInt(256)
	gCoef, _ := GenerateRandInt(256)
	bCoef, _ := GenerateRandInt(256)
	r := uint8(rCoef)
	g := uint8(gCoef)
	b := uint8(bCoef)

	return color.RGBA{R: r, G: g, B: b, A: 255}
}
