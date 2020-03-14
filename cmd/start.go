package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tkiraly/nmea/gpgga"
	"github.com/tkiraly/nmea/gprmc"
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
		f, err := os.OpenFile(pipeFile, os.O_RDWR|os.O_APPEND, 0777)
		if err != nil {
			return err
		}
		for {
			select {
			case <-time.After(time.Second):
				now := time.Now()
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
			case <-s:
				return nil
			}
		}
	},
}
