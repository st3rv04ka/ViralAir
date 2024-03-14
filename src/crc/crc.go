package crc

import "adsb/v2/src/misc"

// For all Mode-S messages

const (
	GENERATOR = "1111111111111010000001001"
)

func Crc(message string, encode bool) string {
	binaryMessageArray := []rune(misc.Hexadecimal2Binary(message))
	generator := []int{}

	for _, char := range GENERATOR {
		generator = append(generator, int(char-'0'))
	}

	if encode {
		for i := len(binaryMessageArray) - 24; i < len(binaryMessageArray); i++ {
			binaryMessageArray[i] = '0'
		}
	}

	for i := 0; i < len(binaryMessageArray)-24; i++ {
		if binaryMessageArray[i] == '1' {
			for j := range generator {
				binaryMessageArray[i+j] = rune('0' + (int(binaryMessageArray[i+j]-'0') ^ generator[j]))
			}
		}
	}

	reminder := string(binaryMessageArray[len(binaryMessageArray)-24:])
	return reminder
}
