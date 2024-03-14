package modulator

import (
	"adsb/v2/src/misc"
)

func Frame1090esPpmModulateCPR(ADSBEvenFrame, ADSBOddFrame []byte) []byte {
	var PulsePositionMululationArray []byte

	for i := 0; i < 48; i++ {
		PulsePositionMululationArray = append(PulsePositionMululationArray, 0)
	}

	PulsePositionMululationArray = append(PulsePositionMululationArray, 0xA1, 0x40)

	for _, byteVal := range ADSBEvenFrame {
		word16 := misc.PackBits(manchesterEncode(^byteVal))
		PulsePositionMululationArray = append(PulsePositionMululationArray, word16[0])
		PulsePositionMululationArray = append(PulsePositionMululationArray, word16[1])
	}

	for i := 0; i < 100; i++ {
		PulsePositionMululationArray = append(PulsePositionMululationArray, 0)
	}

	PulsePositionMululationArray = append(PulsePositionMululationArray, 0xA1, 0x40)

	for _, byteVal := range ADSBOddFrame {
		word16 := misc.PackBits(manchesterEncode(^byteVal))
		PulsePositionMululationArray = append(PulsePositionMululationArray, word16[0])
		PulsePositionMululationArray = append(PulsePositionMululationArray, word16[1])
	}

	for i := 0; i < 48; i++ {
		PulsePositionMululationArray = append(PulsePositionMululationArray, 0)
	}

	return PulsePositionMululationArray
}

func PulsePositionMululation(ADSBFrame []byte) []byte {
	var PulsePositionMululationArray []byte

	for i := 0; i < 48; i++ {
		PulsePositionMululationArray = append(PulsePositionMululationArray, 0)
	}

	PulsePositionMululationArray = append(PulsePositionMululationArray, 0xA1, 0x40)

	for _, byteVal := range ADSBFrame {
		word16 := misc.PackBits(manchesterEncode(^byteVal))
		PulsePositionMululationArray = append(PulsePositionMululationArray, word16[0])
		PulsePositionMululationArray = append(PulsePositionMululationArray, word16[1])
	}

	for i := 0; i < 100; i++ {
		PulsePositionMululationArray = append(PulsePositionMululationArray, 0)
	}

	return PulsePositionMululationArray
}

func manchesterEncode(byte byte) []int {
	var manchesterArray []int

	for i := 7; i >= 0; i-- {
		if misc.ExtractBit(byte, i) {
			manchesterArray = append(manchesterArray, 0, 1)
		} else {
			manchesterArray = append(manchesterArray, 1, 0)
		}
	}

	return manchesterArray
}

func GenerateSDROutput(PulsePositionMululationArray []byte) []byte {
	PulsePositionMululationBits := misc.UnpackBits(PulsePositionMululationArray)
	var SDROutputSignal []byte

	for _, bit := range PulsePositionMululationBits {
		var I, Q byte
		if bit == 1 {
			I = byte(127)
			Q = byte(127)
		} else {
			I = 0
			Q = 0
		}
		SDROutputSignal = append(SDROutputSignal, I, Q)
	}

	return SDROutputSignal
}
