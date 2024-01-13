package modulator

import (
	"adsb/v2/src/misc"
)

func Frame1090esPpmModulate(even, odd []byte) []byte {
	var ppm []byte

	for i := 0; i < 48; i++ {
		ppm = append(ppm, 0)
	}

	ppm = append(ppm, 0xA1, 0x40)

	for _, byteVal := range even {
		word16 := misc.Packbits(manchesterEncode(^byteVal))
		ppm = append(ppm, word16[0])
		ppm = append(ppm, word16[1])
	}

	for i := 0; i < 100; i++ {
		ppm = append(ppm, 0)
	}

	ppm = append(ppm, 0xA1, 0x40)

	for _, byteVal := range odd {
		word16 := misc.Packbits(manchesterEncode(^byteVal))
		ppm = append(ppm, word16[0])
		ppm = append(ppm, word16[1])
	}

	for i := 0; i < 48; i++ {
		ppm = append(ppm, 0)
	}

	return ppm
}

func manchesterEncode(byte byte) []int {
	var manchesterEncoded []int

	for i := 7; i >= 0; i-- {
		if misc.ExtractBit(byte, i) {
			manchesterEncoded = append(manchesterEncoded, 0, 1)
		} else {
			manchesterEncoded = append(manchesterEncoded, 1, 0)
		}
	}

	return manchesterEncoded
}

func GenerateSDROutput(ppm []byte) []byte {
	bits := misc.Unpackbits(ppm)
	var signal []byte

	for _, bit := range bits {
		var I, Q byte
		if bit == 1 {
			I = byte(127)
			Q = byte(127)
		} else {
			I = 0
			Q = 0
		}
		signal = append(signal, I, Q)
	}

	return signal
}
