package pkg

import (
	"crypto/rand"
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
