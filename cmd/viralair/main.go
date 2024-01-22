package main

import (
	"adsb/v2/src/adsb"
	"adsb/v2/src/misc"
	"adsb/v2/src/modulator"
	"flag"
	"fmt"
	"os"
)

func main() {
	// Mode selector
	if len(os.Args) < 2 {
		fmt.Println("Usage: ./main <mode>")
		fmt.Println("Modes:")
		fmt.Println("- adsb")
		os.Exit(1)
	}
	mode := os.Args[1]

	switch mode {
	case "adsb":
		if len(os.Args) < 3 {
			fmt.Println("Usage: ./main <mode> <type>")
			fmt.Println("Types:")
			fmt.Println("- Airborne Position Message (baro) [apmb]")
			fmt.Println("- Aircraft Identification and Category [aiac]")
			os.Exit(1)
		}
		adsbType := os.Args[2]
		switch adsbType {
		case "aiac":
			aiacFlags := flag.NewFlagSet("aiac", flag.ExitOnError)
			var (
				icao = aiacFlags.Int("icao", 0xDEADBE, "icao for ads-s signal")
				tc   = aiacFlags.Int("tc", 4, "type code")
				/**
				1 - Category Set D
				2 - Category Set C
				3 - Category Set B
				4 - Category Set A
				**/
				ca  = aiacFlags.Int("ca", 5, "transponder capability class")
				cat = aiacFlags.Int("cat", 5, "aircraft category")
				/**
				0 - No Information Provided
				1 - Light (less than 7000 kg)
				2 - Medium 1 (between 7000 kg and 34000 kg)
				3 - Medium 2 (between 34000 kg to 136000 kg)
				5 - Heavy (larger than 136000 kg)
				**/
				sign = aiacFlags.String("sign", "XXX777", "aircraft identification 8 chars (@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_ !\"#$%&'()*+,-./0123456789:;<=>?)")
			)
			aiacFlags.Parse(os.Args[3:])

			signEncoded := adsb.GetIdentificationMessage(
				*icao,
				*tc,
				*ca,
				*sign,
				*cat,
			)
			frame := modulator.Frame1090esPpmModulate(signEncoded)
			sdrOutput := modulator.GenerateSDROutput(frame)
			misc.SaveToFile("Samples.iq8s", sdrOutput)
			return

		case "apmb":
			adsbFlags := flag.NewFlagSet("adsb", flag.ExitOnError)
			var (
				icao = adsbFlags.Int("icao", 0xDEADBE, "icao for ads-s signal")
				lat  = adsbFlags.Float64("latitude", 11.33, "aricraft latitude")
				lon  = adsbFlags.Float64("longitude", 11.22, "aircraft longitude")
				alt  = adsbFlags.Float64("altitude", 9999.0, "aircraft altitude")
				ca   = adsbFlags.Int("ca", 5, "transponder capability class")
				tc   = adsbFlags.Int("tc", 11, "type code")
				// The Horizontal Containment Limit (RC) is a parameter that defines the radius
				// of a horizontal area around the true position of an aircraft, within which
				// the indicated position is statistically guaranteed to be contained with a high
				// level of confidence. In aviation navigation, RC provides pilots and air traffic
				// controllers with an indication of the reliability of GPS data.
				/**
				Message format types (tc)
				9  - RC < 7.5 m
				10 - RC < 75 m
				11 - RC < 0.1 NM (185.2 m)
				12 - RC < 0.2 NM (370.4 m)
				13 - RC < 0.6 NM (1111.2 m)
				14 - RC < 1.0 NM (1852 m)
				15 - RC < 2 NM (3.704 km)
				16 - RC < 8 NM (14.816 km)
				17 - RC < 20 NM (37.04 km)
				18 - RC = 20 NM (37.04 km) or unknown
				**/
				ss      = adsbFlags.Int("ss", 0, "surveillance status")
				nicsb   = adsbFlags.Int("nicsb", 0, "navigation integrity category subfield")
				time    = adsbFlags.Int("time", 0, "just ads-b time")
				surface = adsbFlags.Int("surface", 0, "airctaft position (ground/air) (1/0)")
			)
			adsbFlags.Parse(os.Args[3:])

			// ADS-B
			// Get encoded aircraft position
			even, odd := adsb.GetEncodedPosition(
				*ca,
				*icao,
				*tc,
				*ss,
				*nicsb,
				*alt,
				*time,
				*lat,
				*lon,
				*surface,
			)
			frames := modulator.Frame1090esPpmModulateCPR(even, odd)
			sdrOutput := modulator.GenerateSDROutput(frames)
			misc.SaveToFile("Samples.iq8s", sdrOutput)
			return

		}
	default:
		fmt.Println("Usage: ./main <mode>")
		fmt.Println("Modes:")
		fmt.Println("- adsb")
	}

}
