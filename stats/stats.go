package stats

import (
	"math"

	"golang.org/x/exp/constraints"
)

type Numbers interface {
	constraints.Integer | constraints.Float
}

func ArithmeticMean[T Numbers](numbers ...T) float64 {
	if len(numbers) == 0 {
		return 0.0
	}
	sum := 0.0
	for _, number := range numbers {
		sum += float64(number)
	}
	return sum / float64(len(numbers))
}

func StandardDeviation[T Numbers](numbers ...T) float64 {
	if len(numbers) == 0 {
		return 0.0
	}
	mean := ArithmeticMean(numbers...)
	sum := 0.0
	for _, number := range numbers {
		sum += math.Pow(float64(number)-mean, 2)
	}
	return math.Sqrt(sum / float64(len(numbers)))
}
