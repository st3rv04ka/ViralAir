package cpr

import "math"

const (
	latz = 15
)

func nz(ctype int) int {
	return 4*latz - ctype
}

func encodeLatintude(ctype int, surface int) float64 {
	var tmp float64
	if surface == 1 {
		tmp = 90.0
	} else {
		tmp = 360.0
	}

	nzcalc := nz(ctype)
	if nzcalc == 0 {
		return tmp
	} else {
		return tmp / float64(nzcalc)
	}
}

func nl(declatIn float64) float64 {
	if math.Abs(declatIn) >= 87.0 {
		return 1.0
	}
	return math.Floor(
		(2.0 * math.Pi) * math.Pow(math.Acos(1.0-(1.0-math.Cos(math.Pi/(2.0*latz)))/math.Pow(math.Cos((math.Pi/180.0)*math.Abs(declatIn)), 2)), -1))
}

func encodeLongitude(latintude float64, ctype int, surface int) float64 {
	var tmp float64
	if surface == 1 {
		tmp = 90.0
	} else {
		tmp = 360.0
	}
	nlcalc := math.Max(nl(latintude)-float64(ctype), 1)
	return tmp / nlcalc
}

// Encode CPR
func CPREncode(lat float64, lon float64, ctype int, surface int) (int, int) {
	var scalar float64
	if surface == 1 {
		scalar = math.Pow(2, 19)
	} else {
		scalar = math.Pow(2, 17)
	}

	latintude := encodeLatintude(ctype, surface)
	latintudeZonesCount := math.Floor(scalar*((math.Mod(lat, latintude))/latintude) + 0.5)

	longitude := encodeLongitude(lat, ctype, surface)
	longitudeZonesCount := math.Floor(scalar*((math.Mod(lon, longitude))/longitude) + 0.5)

	return int(latintudeZonesCount) & ((1 << 17) - 1), int(longitudeZonesCount) & ((1 << 17) - 1)
}
