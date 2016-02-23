package common

type Size [2]float32

func (s Size) Width() float32 {
	return s[0]
}

func (s Size) Height() float32 {
	return s[1]
}
