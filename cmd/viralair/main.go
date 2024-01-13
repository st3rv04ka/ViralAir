package main

import (
	"adsb/v2/src/adsb"
	"adsb/v2/src/misc"
	"adsb/v2/src/modulator"
	"flag"
)

func main() {
	// Init cmd args
	var (
		icao    = flag.Int("icao", 0xDEADBE, "icao for ads-s signal")
		lat     = flag.Float64("latitude", 11.33, "aricraft latitude")
		lon     = flag.Float64("longitude", 11.22, "aircraft longitude")
		alt     = flag.Float64("altitude", 9999.0, "aircraft altitude")
		ca      = flag.Int("ca", 5, "transponder capability class")
		tc      = flag.Int("tc", 11, "type code")
		ss      = flag.Int("ss", 0, "surveillance status")
		nicsb   = flag.Int("nicsb", 0, "navigation integrity category subfield")
		time    = flag.Int("time", 0, "just ads-b time")
		surface = flag.Int("surface", 0, "airctaft position (ground/air) (1/0)")
		format  = flag.Int("format", 17, "message format. Supported: 17 (adsb)")
	)
	flag.Parse()

	switch *format {
	case 17:
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
		frames := modulator.Frame1090esPpmModulate(even, odd)
		sdrOutput := modulator.GenerateSDROutput(frames)
		misc.SaveToFile("Samples.iq8s", sdrOutput)
		return
	}
}
