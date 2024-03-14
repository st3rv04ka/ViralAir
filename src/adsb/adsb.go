package adsb

import (
	"adsb/v2/src/cpr"
	"adsb/v2/src/crc"
	"adsb/v2/src/misc"
	"fmt"
	"log"
	"strings"
)

const (
	// ADSB Data Format (17)
	format = 17
)

// Encode altitude
func encodeAltitudeModes(alt float64, surface int) int {
	mbit := 0
	qbit := 1
	encalt := int((int(alt) + 1000) / 25)

	var tmp1, tmp2 int

	if surface == 1 {
		tmp1 = (encalt & 0xfe0) << 2
		tmp2 = (encalt & 0x010) << 1
	} else {
		tmp1 = (encalt & 0xff8) << 1
		tmp2 = 0
	}

	return (encalt & 0x0F) | tmp1 | tmp2 | (mbit << 6) | (qbit << 4)
}

// Encode aircraft idenrification message
func GetIdentificationMessage(
	icao int,
	tc int,
	ca int,
	sign string,
	cat int,
) (signEncodedArray []byte) {

	if len(sign) > 8 {
		log.Println("[-] Sign must be less than 8 chars")
		return
	}

	if len(sign) < 8 {
		sign += strings.Repeat(" ", 8-len(sign))
	}

	aicCharset := "@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_ !\"#$%&'()*+,-./0123456789:;<=>?"

	// Format + CA + ICAO
	signEncodedArray = append(signEncodedArray, byte((format<<3)|ca))
	signEncodedArray = append(signEncodedArray, byte((icao>>16)&0xff))
	signEncodedArray = append(signEncodedArray, byte((icao>>8)&0xff))
	signEncodedArray = append(signEncodedArray, byte((icao)&0xff))

	// TC + CAT
	signEncodedArray = append(signEncodedArray, byte((tc<<3)|(cat)))

	// SIGN
	symbols := make([]int, 0, 8)
	for i := 0; i < 8; i++ {
		charPosition := strings.IndexByte(aicCharset, sign[i])
		log.Printf("[+] Encoded char %d -> %02x", sign[i], charPosition)
		symbols = append(symbols, charPosition)
	}
	signEncodedArray = append(signEncodedArray, byte((symbols[0]<<2)|(symbols[1]>>4)))
	signEncodedArray = append(signEncodedArray, byte((symbols[1]<<4)|(symbols[2]>>2)))
	signEncodedArray = append(signEncodedArray, byte((symbols[2]<<6)|symbols[3]))
	signEncodedArray = append(signEncodedArray, byte((symbols[4]<<2)|(symbols[5]>>4)))
	signEncodedArray = append(signEncodedArray, byte((symbols[5]<<4)|(symbols[6]>>2)))
	signEncodedArray = append(signEncodedArray, byte((symbols[6]<<6)|symbols[7]))

	// Convert to hex
	var sbOdd strings.Builder
	for _, b := range signEncodedArray {
		sbOdd.WriteString(fmt.Sprintf("%02x", b))
	}
	signString := sbOdd.String()
	log.Printf("[+] Sign frame without CRC [%s]", signString)

	signCRC := misc.Binary2Integer(crc.Crc(signString+"000000", true))
	log.Printf("[+] Sign frame CRC [%02x]", signCRC)

	signEncodedArray = append(signEncodedArray, byte((signCRC>>16)&0xff))
	signEncodedArray = append(signEncodedArray, byte((signCRC>>8)&0xff))
	signEncodedArray = append(signEncodedArray, byte((signCRC)&0xff))
	log.Printf("[+] Sign frame data [%02x]", signEncodedArray)

	return
}

/**

Based on https://github.com/lyusupov/ADSB-Out

**/

// Encode aircraft position with CPR
func GetEncodedPosition(
	ca int,
	icao int,
	tc int,
	ss int,
	nicsb int,
	alt float64,
	time int,
	lat float64,
	lon float64,
	surface int,
) ([]byte, []byte) {

	// Altitude
	log.Printf("[+] Encode alltitude [%f] with the surface flag [%d]", alt, surface)
	encAlt := encodeAltitudeModes(alt, surface)
	log.Printf("[+] Encoded altirude [0x%02x]", encAlt)

	// Posistion
	// Even
	log.Printf("[+] Encode even frame with lat [%f] and lon [%f]", lat, lon)
	evenLat, evenLon := cpr.CPREncode(lat, lon, 0, surface)
	log.Printf("[+] Encoded even frame lat [0x%02x] and lon [0x%02x]", evenLat, evenLon)
	// Odd
	log.Printf("[+] Encode odd frame with lat [%f] and lon [%f]", lat, lon)
	oddLat, oddLon := cpr.CPREncode(lat, lon, 1, surface)
	log.Printf("[+] Encoded odd frame lat [0x%02x] and lon [0x%02x]", oddLat, oddLon)

	// Encode even data
	ff := 0
	var dataEvenArray []byte
	// Format + CA + ICAO
	dataEvenArray = append(dataEvenArray, byte((format<<3)|ca))
	dataEvenArray = append(dataEvenArray, byte((icao>>16)&0xff))
	dataEvenArray = append(dataEvenArray, byte((icao>>8)&0xff))
	dataEvenArray = append(dataEvenArray, byte((icao)&0xff))

	// Lat + Lot + Alt (even)
	dataEvenArray = append(dataEvenArray, byte((tc<<3)|(ss<<1)|nicsb))
	dataEvenArray = append(dataEvenArray, byte((encAlt>>4)&0xff))
	dataEvenArray = append(dataEvenArray, byte((encAlt&0xf)<<4|(time<<3)|(ff<<2)|(evenLat>>15)))
	dataEvenArray = append(dataEvenArray, byte((evenLat>>7)&0xff))
	dataEvenArray = append(dataEvenArray, byte(((evenLat&0x7f)<<1)|(evenLon>>16)))
	dataEvenArray = append(dataEvenArray, byte((evenLon>>8)&0xff))
	dataEvenArray = append(dataEvenArray, byte((evenLon)&0xff))

	// Convert to hex
	var sbEven strings.Builder
	for _, b := range dataEvenArray[:11] {
		sbEven.WriteString(fmt.Sprintf("%02x", b))
	}
	dataEvenString := sbEven.String()
	log.Printf("[+] Even frame without CRC [%s]", dataEvenString)

	// Calculate CRC
	dataEvenCRC := misc.Binary2Integer(crc.Crc(dataEvenString+"000000", true))
	log.Printf("[+] Even data CRC [%02x]", dataEvenCRC)

	// Append CRC
	dataEvenArray = append(dataEvenArray, byte((dataEvenCRC>>16)&0xff))
	dataEvenArray = append(dataEvenArray, byte((dataEvenCRC>>8)&0xff))
	dataEvenArray = append(dataEvenArray, byte((dataEvenCRC)&0xff))
	log.Printf("[+] Even data [%02x]", dataEvenArray)

	// Encode odd data
	ff = 1
	var dataOddArray []byte
	// Format + CA + ICAO
	dataOddArray = append(dataOddArray, byte((format<<3)|ca))
	dataOddArray = append(dataOddArray, byte((icao>>16)&0xff))
	dataOddArray = append(dataOddArray, byte((icao>>8)&0xff))
	dataOddArray = append(dataOddArray, byte((icao)&0xff))

	// Lat + Lot + Alt (even)
	dataOddArray = append(dataOddArray, byte((tc<<3)|(ss<<1)|nicsb))
	dataOddArray = append(dataOddArray, byte((encAlt>>4)&0xff))
	dataOddArray = append(dataOddArray, byte((encAlt&0xf)<<4|(time<<3)|(ff<<2)|(oddLat>>15)))
	dataOddArray = append(dataOddArray, byte((oddLat>>7)&0xff))
	dataOddArray = append(dataOddArray, byte(((oddLat&0x7f)<<1)|(oddLon>>16)))
	dataOddArray = append(dataOddArray, byte((oddLon>>8)&0xff))
	dataOddArray = append(dataOddArray, byte((oddLon)&0xff))

	// Convert to hex
	var sbOdd strings.Builder
	for _, b := range dataOddArray[:11] {
		sbOdd.WriteString(fmt.Sprintf("%02x", b))
	}
	dataOddString := sbOdd.String()
	log.Printf("[+] Odd frame without CRC [%s]", dataOddString)

	dataOddCRC := misc.Binary2Integer(crc.Crc(dataOddString+"000000", true))
	log.Printf("[+] Odd data CRC [%02x]", dataOddCRC)

	dataOddArray = append(dataOddArray, byte((dataOddCRC>>16)&0xff))
	dataOddArray = append(dataOddArray, byte((dataOddCRC>>8)&0xff))
	dataOddArray = append(dataOddArray, byte((dataOddCRC)&0xff))
	log.Printf("[+] Odd data [%02x]", dataOddArray)

	return dataEvenArray, dataOddArray
}
