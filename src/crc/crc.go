package crc

import "adsb/v2/src/misc"

// For all Mode-S messages

const (
	GENERATOR = "1111111111111010000001001"
)

func Crc(msg string, encode bool) string {
	msgbin := []rune(misc.Hex2bin(msg))
	generator := []int{}

	for _, char := range GENERATOR {
		generator = append(generator, int(char-'0'))
	}

	if encode {
		for i := len(msgbin) - 24; i < len(msgbin); i++ {
			msgbin[i] = '0'
		}
	}

	for i := 0; i < len(msgbin)-24; i++ {
		if msgbin[i] == '1' {
			for j := range generator {
				msgbin[i+j] = rune('0' + (int(msgbin[i+j]-'0') ^ generator[j]))
			}
		}
	}

	reminder := string(msgbin[len(msgbin)-24:])
	return reminder
}
