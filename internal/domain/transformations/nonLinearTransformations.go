package transformations

import "math"

func Spherical(x, y float64) (newX, newY float64) {
	r := x*x + y*y
	if r == 0 {
		return 0, 0 // Защита от деления на 0
	}

	return x / r, y / r
}

func Sinusoidal(x, y float64) (newX, newY float64) {
	newX = math.Sin(x * math.Pi)
	newY = math.Sin(y * math.Pi)

	return
}

func Handkerchief(x, y float64) (newX, newY float64) {
	r := math.Sqrt((x * x) + (y * y))
	theta := math.Atan2(y, x)

	newX = r * math.Sin(theta+r)
	newY = math.Cos(theta - r)

	return
}

func Swirl(x, y float64) (newX, newY float64) {
	r := math.Sqrt(x*x + y*y)
	newX = x*math.Sin(r*r) - y*math.Cos(r*r)
	newY = x*math.Cos(r*r) - y*math.Sin(r*r)

	return
}

func Horseshoe(x, y float64) (newX, newY float64) {
	r := math.Sqrt(x*x + y*y)
	if r == 0 {
		return 0, 0
	}

	newX = (x - y) * (x + y) / r
	newY = 2 * x * y / r

	return
}

func Polar(x, y float64) (newX, newY float64) {
	r := math.Sqrt(x*x + y*y)
	theta := math.Atan2(y, x)
	newX = theta / math.Pi
	newY = r - 1

	return
}

func Disc(x, y float64) (newX, newY float64) {
	r := math.Sqrt(x*x + y*y)
	theta := math.Atan2(y, x)
	newX = theta / math.Pi * math.Sin(math.Pi*r)
	newY = math.Cos(math.Pi * r)

	return
}

func Heart(x, y float64) (newX, newY float64) {
	r := math.Sqrt(x*x + y*y)
	theta := math.Atan2(y, x)
	newX = r * math.Sin(theta*r)
	newY = -r * math.Cos(theta*r)

	return
}

func Linear(x, y float64) (newX, newY float64) {
	return x, y
}

func EyeFish(x, y float64) (newX, newY float64) {
	r := math.Sqrt(x*x + y*y)
	newX = 2.0 / (r + 1) * x
	newY = 2.0 / (r + 1) * y

	return
}
