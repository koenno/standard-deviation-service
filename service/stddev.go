package service

import (
	"github.com/koenno/standard-deviation-service/stats"
)

type StdDevService struct {
}

func NewStdDevService() StdDevService {
	return StdDevService{}
}

type StdDevResult struct {
	StdDev float64 `json:"stddev"`
	Data   []int   `json:"data"`
}

func (s StdDevService) Calculate(input <-chan []int) <-chan StdDevResult {
	output := make(chan StdDevResult)
	go func() {
		defer close(output)
		var setSum []int
		for set := range input {
			setSum = append(setSum, set...)
			stddev := stats.StandardDeviation(set...)
			output <- StdDevResult{
				StdDev: stddev,
				Data:   set,
			}
		}
		if len(setSum) == 0 {
			return
		}
		stddev := stats.StandardDeviation(setSum...)
		output <- StdDevResult{
			StdDev: stddev,
			Data:   setSum,
		}
	}()
	return output
}
