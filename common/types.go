package common

import "time"

type Time int64

func Now() Time {
	return Time(time.Now().UTC().UnixNano())
}

func (t Time) Seconds() float64 {
	return float64(t) / 1000000000
}

type Size [2]float32

func (s Size) Width() float32 {
	return s[0]
}

func (s Size) Height() float32 {
	return s[1]
}
