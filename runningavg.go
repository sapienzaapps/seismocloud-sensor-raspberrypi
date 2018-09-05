package main

import "math"

type RunningAvgFloat64 interface {
	AddValue(float64) float64
	GetAverage() float64
	GetVariance() float64
	GetStandardDeviation() float64
	Elements() uint64
}

func NewRunningAvgFloat64() RunningAvgFloat64 {
	return &runningAvgFloat64impl{}
}

type runningAvgFloat64impl struct {
	avg      float64
	variance float64
	n        uint64
}

func (r *runningAvgFloat64impl) AddValue(value float64) float64 {
	r.n += 1
	delta := value - r.avg
	r.avg += delta / float64(r.n)
	r.variance = r.variance + delta*(value-r.avg)
	return r.avg
}

func (r *runningAvgFloat64impl) GetAverage() float64 {
	return r.avg
}

func (r *runningAvgFloat64impl) GetVariance() float64 {
	if r.n < 2 {
		return 0
	}
	return r.variance / float64(r.n-1)
}

func (r *runningAvgFloat64impl) GetStandardDeviation() float64 {
	return math.Sqrt(r.GetVariance())
}

func (r *runningAvgFloat64impl) Elements() uint64 {
	return r.n
}
