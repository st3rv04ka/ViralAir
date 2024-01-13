package misc

import (
	"fmt"
	"math"
	"math/big"
	"os"
	"strings"
)

func Bin2int(binStr string) int64 {
	n := new(big.Int)
	n, ok := n.SetString(binStr, 2)
	if !ok {
		return 0
	}
	return n.Int64()
}

func ExtractBit(byte byte, position int) bool {
	return (byte>>position)&1 == 1
}

// Numpy pack
func Packbits(bits []int) []byte {
	numBytes := int(math.Ceil(float64(len(bits)) / 8.0))
	packedBytes := make([]byte, numBytes)

	for i, bit := range bits {
		if bit != 0 {
			byteIndex := i / 8
			bitPosition := uint(7 - i%8)
			packedBytes[byteIndex] |= 1 << bitPosition
		}
	}

	return packedBytes
}

func SaveToFile(filename string, data []byte) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}

// Numpy unpack
func Unpackbits(bytes []byte) []int {
	var bits []int
	for _, b := range bytes {
		for i := 7; i >= 0; i-- {
			bits = append(bits, int((b>>uint(i))&1))
		}
	}
	return bits
}

func Hex2bin(hexStr string) string {
	n := new(big.Int)
	n, ok := n.SetString(hexStr, 16)
	if !ok {
		return ""
	}

	binStr := fmt.Sprintf("%b", n)

	numBits := len(hexStr) * 4
	if len(binStr) < numBits {
		binStr = strings.Repeat("0", numBits-len(binStr)) + binStr
	}

	return binStr
}
