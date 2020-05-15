package gprmc

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/tkiraly/fakegps/nmea/util"
)

type GPRMC struct {
	DateTime          time.Time
	ReceiverStatus    string // A = OK, V = Warning
	Latitude          float64
	Longitude         float64
	Speed             float64
	Course            float64
	MagneticVariation float64
}

func Build(sentence GPRMC) (string, error) {
	lat, lon := util.Lonlat2nmealonlat(sentence.Longitude, sentence.Latitude)
	magvardir := "A"
	if sentence.MagneticVariation < 0 {
		magvardir = "W"
	}
	s := fmt.Sprintf("GPRMC,%02d%02d%02d.%03d,%s,%s,%s,%05.1f,%05.1f,%02d%02d%02d,%05.1f,,%s",
		sentence.DateTime.UTC().Hour(), sentence.DateTime.UTC().Minute(), sentence.DateTime.UTC().Second(),
		sentence.DateTime.UTC().Nanosecond()/1_000_000,
		sentence.ReceiverStatus, lat, lon, sentence.Speed, sentence.Course,
		sentence.DateTime.Day(), sentence.DateTime.Month(), sentence.DateTime.Year()%100,
		sentence.MagneticVariation, magvardir)
	s = util.Checksum(s)
	return s, nil
}

func BuildMinimal(now time.Time, longitude, latitude float64) (string, error) {
	rmc := GPRMC{
		DateTime:          now,
		Longitude:         longitude,
		Latitude:          latitude,
		ReceiverStatus:    "A",
		Speed:             0, // Speed in knots
		Course:            0,
		MagneticVariation: 0,
	}
	return Build(rmc)
}

func Parse(sentence string, v *GPRMC) error {
	components := strings.Split(sentence, ",")
	t, err := util.Parsedatetime(components[1], components[9])
	if err != nil {
		return err
	}

	latint, err := strconv.Atoi(components[3][:2])
	if err != nil {
		return err
	}
	latfrac, err := strconv.ParseFloat(components[3][2:], 64)
	latfrac = latfrac / 60.0
	lat := float64(latint) + latfrac
	if components[4] == "S" {
		lat = -lat
	}
	lonint, err := strconv.Atoi(components[5][:3])
	if err != nil {
		return err
	}
	lonfrac, err := strconv.ParseFloat(components[5][3:], 64)
	lonfrac = lonfrac / 60.0
	lon := float64(lonint) + lonfrac
	if components[6] == "W" {
		lon = -lon
	}
	speed, err := util.ParseFloat(components[7])
	if err != nil {
		return err
	}
	course, err := util.ParseFloat(components[8])
	if err != nil {
		return err
	}
	magvar, err := util.ParseFloat(components[10])
	if err != nil {
		return err
	}
	if components[11] == "W" {
		magvar = -magvar
	}
	indexofstar := strings.Index(sentence, "*")
	payload := sentence[1:indexofstar]
	if util.Checksum(payload) != sentence {
		return fmt.Errorf("checksum error")
	}
	v.DateTime = t
	v.ReceiverStatus = components[2]
	v.Latitude = lat
	v.Longitude = lon
	v.Speed = speed
	v.Course = course
	v.MagneticVariation = magvar
	return nil
}
