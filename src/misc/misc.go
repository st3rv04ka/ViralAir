package misc

import (
	"fmt"
	"math"
	"math/big"
	"os"
	"strings"
)

func Binary2Integer(binary string) int64 {
	integer := new(big.Int)
	integer, ok := integer.SetString(binary, 2)
	if !ok {
		return 0
	}
	return integer.Int64()
}

func ExtractBit(byte byte, bitPosition int) bool {
	return (byte>>bitPosition)&1 == 1
}

// Numpy pack
func PackBits(bits []int) []byte {
	bytesCount := int(math.Ceil(float64(len(bits)) / 8.0))
	packedBytesArray := make([]byte, bytesCount)

	for position, bit := range bits {
		if bit != 0 {
			byteIndex := position / 8
			bitPosition := uint(7 - position%8)
			packedBytesArray[byteIndex] |= 1 << bitPosition
		}
	}

	return packedBytesArray
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
func UnpackBits(bytes []byte) []int {
	var bitsArray []int
	for _, b := range bytes {
		for i := 7; i >= 0; i-- {
			bitsArray = append(bitsArray, int((b>>uint(i))&1))
		}
	}
	return bitsArray
}

func Hexadecimal2Binary(hex string) string {
	binary := new(big.Int)
	binary, ok := binary.SetString(hex, 16)
	if !ok {
		return ""
	}

	binaryString := fmt.Sprintf("%b", binary)

	bitsCount := len(hex) * 4
	if len(binaryString) < bitsCount {
		binaryString = strings.Repeat("0", bitsCount-len(binaryString)) + binaryString
	}

	return binaryString
}
