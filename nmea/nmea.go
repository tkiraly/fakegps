package nmea

import (
	"fmt"

	"github.com/tkiraly/fakegps/nmea/gpgga"
	"github.com/tkiraly/fakegps/nmea/gprmc"
)

func Parse(sentence string, v interface{}) error {
	switch t := v.(type) {
	case *gpgga.GPGGA:
		return gpgga.Parse(sentence, t)
	case *gprmc.GPRMC:
		return gprmc.Parse(sentence, t)
	default:
		return fmt.Errorf("unknown sentence type: %T", t)
	}
}

func Build(sentence interface{}) (string, error) {
	switch v := sentence.(type) {
	case gpgga.GPGGA:
		return gpgga.Build(v)
	case *gpgga.GPGGA:
		return gpgga.Build(*v)
	case gprmc.GPRMC:
		return gprmc.Build(v)
	case *gprmc.GPRMC:
		return gprmc.Build(*v)
	default:
		return "", fmt.Errorf("unknown sentence type: %T", sentence)
	}
}
