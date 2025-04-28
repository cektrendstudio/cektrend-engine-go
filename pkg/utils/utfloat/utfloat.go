package utfloat

import "math"

func RoundEx(val float64, dec int32) float64 {
	pow := math.Pow(10, float64(dec))
	dgt := pow * val
	_, div := math.Modf(dgt)

	var round float64
	if val > 0 {
		if div >= 0.5 {
			round = math.Ceil(dgt)
		} else {
			round = math.Floor(dgt)
		}
	} else {
		if div >= 0.5 {
			round = math.Floor(dgt)
		} else {
			round = math.Ceil(dgt)
		}
	}

	return round / pow
}

func Round(val float64, dec int32) float64 {
	addx := float64(1)
	for i := int32(0); i < dec; i++ {
		addx *= 10
	}

	return math.Round(val*addx) / addx
}

func Floor(val float64, dec int32) float64 {
	addx := float64(1)
	for i := int32(0); i < dec; i++ {
		addx *= 10
	}

	return math.Floor(val*addx) / addx
}

func Ceil(val float64, dec int32) float64 {
	addx := float64(1)
	for i := int32(0); i < dec; i++ {
		addx *= 10
	}

	return math.Ceil(val*addx) / addx
}
