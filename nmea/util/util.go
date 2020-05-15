package util

import (
	"fmt"
	"strconv"
	"time"
)

func Parsetime(component string) (time.Time, error) {
	h, err := strconv.Atoi(component[:2])
	if err != nil {
		return time.Time{}, err
	}
	m, err := strconv.Atoi(component[2:4])
	if err != nil {
		return time.Time{}, err
	}
	s, err := strconv.ParseFloat(component[4:], 64)
	if err != nil {
		return time.Time{}, err
	}
	intpart := int(s)
	ratpart := s - float64(intpart)
	return time.Date(2020, 03, 03, h, m, intpart, int(ratpart*1_000_000_000), time.UTC), nil
}
func Parsedatetime(tim, date string) (time.Time, error) {
	h, err := strconv.Atoi(tim[:2])
	if err != nil {
		return time.Time{}, err
	}
	m, err := strconv.Atoi(tim[2:4])
	if err != nil {
		return time.Time{}, err
	}
	s, err := strconv.ParseFloat(tim[4:], 64)
	if err != nil {
		return time.Time{}, err
	}
	day, err := strconv.Atoi(date[:2])
	if err != nil {
		return time.Time{}, err
	}
	month, err := strconv.Atoi(date[2:4])
	if err != nil {
		return time.Time{}, err
	}
	year, err := strconv.Atoi(date[4:])
	if err != nil {
		return time.Time{}, err
	}
	// NMEA protocol was released in 1983 so 83 is a good devider.
	// see you in 2083 :)
	if year < 83 {
		year += 2000
	} else {
		year += 1900
	}
	intpart := int(s)
	ratpart := s - float64(intpart)
	return time.Date(year, time.Month(month), day, h, m, intpart, int(ratpart*1_000_000_000), time.UTC), nil
}

func Lonlat2nmealonlat(lon, lat float64) (string, string) {
	latdir := "N"
	if lat < 0 {
		latdir = "S"
		lat = -lat
	}
	londir := "E"
	if lon < 0 {
		londir = "W"
		lon = -lon
	}
	latmin := (lat - float64(int(lat))) * 60
	lonmin := (lon - float64(int(lon))) * 60
	lati := fmt.Sprintf("%02d%07.04f,%s", int(lat), latmin, latdir)
	long := fmt.Sprintf("%03d%07.04f,%s", int(lon), lonmin, londir)

	return lati, long
}
func Checksum(in string) string {
	checksum := byte(0)
	for i := 0; i < len(in); i++ {
		checksum ^= byte(in[i])
	}
	return fmt.Sprintf("$%s*%X", in, checksum)
}

func ParseFloat(in string) (float64, error) {
	if len(in) == 0 {
		return 0, nil
	}
	n, err := strconv.ParseFloat(in, 64)
	if err != nil {
		return 0, err
	}
	return n, nil
}
func Atoi(in string) (int, error) {
	if len(in) == 0 {
		return 0, nil
	}
	id, err := strconv.Atoi(in)
	if err != nil {
		return 0, err
	}
	return id, nil
}
