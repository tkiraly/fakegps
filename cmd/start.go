package cmd

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tkiraly/fakegps/nmea/gpgga"
	"github.com/tkiraly/fakegps/nmea/gprmc"
)

var gpspath string
var longitude, latitude, hdop, altitude float64
var usednumofsats int

func init() {
	startCmd.Flags().StringVar(&gpspath, "gpspath", "/tmp/fakegps", "path of the fifo to push to")
	viper.BindPFlag("gpspath", startCmd.Flags().Lookup("gpspath"))
	viper.BindEnv("gpspath", "GPS_PATH")
	startCmd.Flags().Float64Var(&longitude, "longitude", 19.078152, "path of the fifo to push to")
	viper.BindPFlag("longitude", startCmd.Flags().Lookup("longitude"))
	viper.BindEnv("longitude", "LONGITUDE")
	startCmd.Flags().Float64Var(&latitude, "latitude", 47.515111, "path of the fifo to push to")
	viper.BindPFlag("latitude", startCmd.Flags().Lookup("latitude"))
	viper.BindEnv("latitude", "LATITUDE")
	startCmd.Flags().Float64Var(&altitude, "altitude", 105, "path of the fifo to push to")
	viper.BindPFlag("altitude", startCmd.Flags().Lookup("altitude"))
	viper.BindEnv("altitude", "ALTITUDE")
	rootCmd.AddCommand(startCmd)
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start simulator",
	RunE: func(cmd *cobra.Command, args []string) error {
		s := make(chan os.Signal, 1)
		signal.Notify(s, os.Interrupt)
		pipeFile := viper.GetString("gpspath")
		_, err := os.Stat(pipeFile)
		if err != nil {
			if os.IsNotExist(err) {
				err := syscall.Mkfifo(pipeFile, 0666)
				if err != nil {
					return err
				}
			} else {
				return err
			}
		} else {
			err := os.Remove(pipeFile)
			if err != nil {
				return err
			}
			err = syscall.Mkfifo(pipeFile, 0666)
			if err != nil {
				return err
			}
		}
		f, err := os.OpenFile(pipeFile, os.O_WRONLY, 0777)
		if err != nil {
			return err
		}
		for {
			select {
			case <-time.After(time.Second):
				now := time.Now().Add(-500 * time.Millisecond)
				fullnmea, err := gpgga.BuildMinimal(now, viper.GetFloat64("longitude"), viper.GetFloat64("latitude"), viper.GetFloat64("altitude"))
				if err != nil {
					return err
				}
				line := fullnmea + "\r\n"
				_, err = f.WriteString(line)
				if err != nil {
					return err
				}
				fmt.Print(line)
				fullnmea, err = gprmc.BuildMinimal(now, viper.GetFloat64("longitude"), viper.GetFloat64("latitude"))
				if err != nil {
					return err
				}
				line = fullnmea + "\r\n"
				_, err = f.WriteString(line)
				if err != nil {
					return err
				}
				fmt.Print(line)
				ubx := [24]byte{
					0xB5, 0x62,
					0x01, 0x20,
					16, 0,
				}
				epoch := time.Date(1980, 1, 6, 0, 0, 0, 0, time.Local)
				dd := now.Sub(epoch)
				week := int16(dd / 1000 / 1000 / 1000 / 60 / 60 / 24 / 7)
				tow := int64(dd % (1000 * 1000 * 1000 * 60 * 60 * 24 * 7))
				ftow := int32(tow % 1_000_000)
				itow := uint32(0)
				if ftow > 500000 {
					//round up itow
					itow = uint32(tow/1000_000) + 1
					ftow = ftow - 1_000_000
				} else {
					itow = uint32(tow / 1000_000)
				}
				binary.LittleEndian.PutUint32(ubx[6:10], itow)
				ubx[10] = byte(ftow)
				ubx[11] = byte(ftow >> 8)
				ubx[12] = byte(ftow >> 16)
				ubx[13] = byte(ftow >> 24)
				binary.LittleEndian.PutUint32(ubx[10:14], uint32(ftow))
				binary.LittleEndian.PutUint16(ubx[14:16], uint16(week))
				ubx[16] = 18
				ubx[17] = 0x07
				binary.LittleEndian.PutUint32(ubx[18:22], 70)
				ck_a := byte(0)
				ck_b := byte(0)

				for i := 0; i < (4 + 16); i++ {
					ck_a = ck_a + ubx[i+2]
					ck_b = ck_b + ck_a
				}
				ubx[22] = ck_a
				ubx[23] = ck_b
				_, err = f.Write(ubx[:])
				if err != nil {
					return err
				}
				f.Sync()
				fmt.Printf("%s\n", hex.EncodeToString(ubx[:]))
			case <-s:
				return nil
			}
		}
	},
}
