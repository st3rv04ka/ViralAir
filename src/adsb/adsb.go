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
func encodeAltModes(alt float64, surface int) int {
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
) (signEncoded []byte) {

	if len(sign) > 8 {
		log.Println("[-] Sign must be less than 8 chars")
		return
	}

	if len(sign) < 8 {
		sign += strings.Repeat(" ", 8-len(sign))
	}

	aicCharset := "@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_ !\"#$%&'()*+,-./0123456789:;<=>?"

	// Format + CA + ICAO
	signEncoded = append(signEncoded, byte((format<<3)|ca))
	signEncoded = append(signEncoded, byte((icao>>16)&0xff))
	signEncoded = append(signEncoded, byte((icao>>8)&0xff))
	signEncoded = append(signEncoded, byte((icao)&0xff))

	// TC + CAT
	signEncoded = append(signEncoded, byte((tc<<3)|(cat)))

	// SIGN
	symbols := make([]int, 0, 8)
	for i := 0; i < 8; i++ {
		charPosition := strings.IndexByte(aicCharset, sign[i])
		log.Printf("[+] Encoded char %d -> %02x", sign[i], charPosition)
		symbols = append(symbols, charPosition)
	}
	signEncoded = append(signEncoded, byte((symbols[0]<<2)|(symbols[1]>>4)))
	signEncoded = append(signEncoded, byte((symbols[1]<<4)|(symbols[2]>>2)))
	signEncoded = append(signEncoded, byte((symbols[2]<<6)|symbols[3]))
	signEncoded = append(signEncoded, byte((symbols[4]<<2)|(symbols[5]>>4)))
	signEncoded = append(signEncoded, byte((symbols[5]<<4)|(symbols[6]>>2)))
	signEncoded = append(signEncoded, byte((symbols[6]<<6)|symbols[7]))

	// Convert to hex
	var sbOdd strings.Builder
	for _, b := range signEncoded {
		sbOdd.WriteString(fmt.Sprintf("%02x", b))
	}
	signString := sbOdd.String()
	log.Printf("[+] Sign frame without CRC [%s]", signString)

	signCRC := misc.Bin2int(crc.Crc(signString+"000000", true))
	log.Printf("[+] Sign frame CRC [%02x]", signCRC)

	signEncoded = append(signEncoded, byte((signCRC>>16)&0xff))
	signEncoded = append(signEncoded, byte((signCRC>>8)&0xff))
	signEncoded = append(signEncoded, byte((signCRC)&0xff))
	log.Printf("[+] Sign frame data [%02x]", signEncoded)

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
) (even []byte, odd []byte) {

	// Altitude
	log.Printf("[+] Encode alltitude [%f] with the surface flag [%d]", alt, surface)
	encAlt := encodeAltModes(alt, surface)
	log.Printf("[+] Encoded altirude [0x%02x]", encAlt)

	// Posistion
	// Even
	log.Printf("[+] Encode even frame with lat [%f] and lon [%f]", lat, lon)
	evenLat, evenLon := cpr.CprEncode(lat, lon, 0, surface)
	log.Printf("[+] Encoded even frame lat [0x%02x] and lon [0x%02x]", evenLat, evenLon)
	// Odd
	log.Printf("[+] Encode odd frame with lat [%f] and lon [%f]", lat, lon)
	oddLat, oddLon := cpr.CprEncode(lat, lon, 1, surface)
	log.Printf("[+] Encoded odd frame lat [0x%02x] and lon [0x%02x]", oddLat, oddLon)

	// Encode even data
	ff := 0
	var dataEven []byte
	// Format + CA + ICAO
	dataEven = append(dataEven, byte((format<<3)|ca))
	dataEven = append(dataEven, byte((icao>>16)&0xff))
	dataEven = append(dataEven, byte((icao>>8)&0xff))
	dataEven = append(dataEven, byte((icao)&0xff))

	// Lat + Lot + Alt (even)
	dataEven = append(dataEven, byte((tc<<3)|(ss<<1)|nicsb))
	dataEven = append(dataEven, byte((encAlt>>4)&0xff))
	dataEven = append(dataEven, byte((encAlt&0xf)<<4|(time<<3)|(ff<<2)|(evenLat>>15)))
	dataEven = append(dataEven, byte((evenLat>>7)&0xff))
	dataEven = append(dataEven, byte(((evenLat&0x7f)<<1)|(evenLon>>16)))
	dataEven = append(dataEven, byte((evenLon>>8)&0xff))
	dataEven = append(dataEven, byte((evenLon)&0xff))

	// Convert to hex
	var sbEven strings.Builder
	for _, b := range dataEven[:11] {
		sbEven.WriteString(fmt.Sprintf("%02x", b))
	}
	dataEvenString := sbEven.String()
	log.Printf("[+] Even frame without CRC [%s]", dataEvenString)

	// Calculate CRC
	dataEvenCRC := misc.Bin2int(crc.Crc(dataEvenString+"000000", true))
	log.Printf("[+] Even data CRC [%02x]", dataEvenCRC)

	// Append CRC
	dataEven = append(dataEven, byte((dataEvenCRC>>16)&0xff))
	dataEven = append(dataEven, byte((dataEvenCRC>>8)&0xff))
	dataEven = append(dataEven, byte((dataEvenCRC)&0xff))
	log.Printf("[+] Even data [%02x]", dataEven)

	// Encode odd data
	ff = 1
	var dataOdd []byte
	// Format + CA + ICAO
	dataOdd = append(dataOdd, byte((format<<3)|ca))
	dataOdd = append(dataOdd, byte((icao>>16)&0xff))
	dataOdd = append(dataOdd, byte((icao>>8)&0xff))
	dataOdd = append(dataOdd, byte((icao)&0xff))

	// Lat + Lot + Alt (even)
	dataOdd = append(dataOdd, byte((tc<<3)|(ss<<1)|nicsb))
	dataOdd = append(dataOdd, byte((encAlt>>4)&0xff))
	dataOdd = append(dataOdd, byte((encAlt&0xf)<<4|(time<<3)|(ff<<2)|(oddLat>>15)))
	dataOdd = append(dataOdd, byte((oddLat>>7)&0xff))
	dataOdd = append(dataOdd, byte(((oddLat&0x7f)<<1)|(oddLon>>16)))
	dataOdd = append(dataOdd, byte((oddLon>>8)&0xff))
	dataOdd = append(dataOdd, byte((oddLon)&0xff))

	// Convert to hex
	var sbOdd strings.Builder
	for _, b := range dataOdd[:11] {
		sbOdd.WriteString(fmt.Sprintf("%02x", b))
	}
	dataOddString := sbOdd.String()
	log.Printf("[+] Odd frame without CRC [%s]", dataOddString)

	dataOddCRC := misc.Bin2int(crc.Crc(dataOddString+"000000", true))
	log.Printf("[+] Odd data CRC [%02x]", dataOddCRC)

	dataOdd = append(dataOdd, byte((dataOddCRC>>16)&0xff))
	dataOdd = append(dataOdd, byte((dataOddCRC>>8)&0xff))
	dataOdd = append(dataOdd, byte((dataOddCRC)&0xff))
	log.Printf("[+] Odd data [%02x]", dataOdd)

	return dataEven, dataOdd
}
