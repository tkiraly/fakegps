package gpgga

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/tkiraly/fakegps/nmea/util"
)

type GPGGA struct {
	Time          time.Time
	Latitude      float64
	Longitude     float64
	FixQuality    int
	NumSatellites int
	HDOP          float64
	Altitude      float64
	Separation    float64
	DGPSAge       float64
	DGPSId        int
}

func Build(sentence GPGGA) (string, error) {
	lat, lon := util.Lonlat2nmealonlat(sentence.Longitude, sentence.Latitude)
	s := fmt.Sprintf("GPGGA,%02d%02d%06.3f,%s,%s,%1d,%02d,%03.1f,%03.1f,%s,%03.1f,%s,%03.1f,%04d",
		sentence.Time.UTC().Hour(), sentence.Time.UTC().Minute(),
		float64(sentence.Time.UTC().Second())+float64(sentence.Time.UTC().Nanosecond())/1_000_000_000,
		lat, lon, sentence.FixQuality, sentence.NumSatellites, sentence.HDOP, sentence.Altitude, "M",
		sentence.Separation, "M", sentence.DGPSAge, sentence.DGPSId)
	s = util.Checksum(s)
	return s, nil
}

func BuildMinimal(now time.Time, longitude, latitude, altitude float64) (string, error) {
	gga := GPGGA{
		Time:          now,
		Longitude:     longitude,
		Latitude:      latitude,
		FixQuality:    1,
		NumSatellites: 10,
		HDOP:          1,
		Altitude:      altitude,
		Separation:    0,
		DGPSAge:       0,
		DGPSId:        0,
	}
	return Build(gga)
}

func Parse(sentence string, v *GPGGA) error {
	components := strings.Split(sentence, ",")
	t, err := util.Parsetime(components[1])
	if err != nil {
		return err
	}
	latint, err := strconv.Atoi(components[2][:2])
	if err != nil {
		return err
	}
	latfrac, err := strconv.ParseFloat(components[2][2:], 64)
	latfrac = latfrac / 60.0
	lat := float64(latint) + latfrac
	if components[3] == "S" {
		lat = -lat
	}
	lonint, err := strconv.Atoi(components[4][:3])
	if err != nil {
		return err
	}
	lonfrac, err := strconv.ParseFloat(components[4][3:], 64)
	lonfrac = lonfrac / 60.0
	lon := float64(lonint) + lonfrac
	if components[5] == "W" {
		lon = -lon
	}
	quality, err := strconv.Atoi(components[6])
	if err != nil {
		return err
	}
	sats, err := strconv.Atoi(components[7])
	if err != nil {
		return err
	}
	hdop, err := strconv.ParseFloat(components[8], 64)
	if err != nil {
		return err
	}
	alt, err := strconv.ParseFloat(components[9], 64)
	if err != nil {
		return err
	}
	sep, err := strconv.ParseFloat(components[11], 64)
	if err != nil {
		return err
	}
	age, err := util.ParseFloat(components[13])
	if err != nil {
		return err
	}
	indexofstar := strings.Index(components[14], "*")
	id, err := util.Atoi(components[14][:indexofstar])
	if err != nil {
		return err
	}
	indexofstar = strings.Index(sentence, "*")
	payload := sentence[1:indexofstar]
	if util.Checksum(payload) != sentence {
		return fmt.Errorf("checksum error")
	}
	v.Time = t
	v.Latitude = lat
	v.Longitude = lon
	v.FixQuality = quality
	v.NumSatellites = sats
	v.HDOP = hdop
	v.Altitude = alt
	v.Separation = sep
	v.DGPSAge = age
	v.DGPSId = id
	return nil
}
