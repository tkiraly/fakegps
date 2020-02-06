package cmd

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var path string

func init() {
	startCmd.Flags().StringVar(&path, "path", "/tmp/fakegps", "path of the fifo to push to")
	viper.BindPFlag("path", startCmd.Flags().Lookup("path"))
	rootCmd.AddCommand(startCmd)
}

func checksum(in string) string {
	checksum := byte(0)
	for i := 0; i < len(in); i++ {
		checksum ^= byte(in[i])
	}
	return fmt.Sprintf("$%s*%X\n", in, checksum)
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start simulator",
	RunE: func(cmd *cobra.Command, args []string) error {
		s := make(chan os.Signal, 1)
		signal.Notify(s, os.Interrupt)
		pipeFile := viper.GetString("path")
		os.Remove(pipeFile)
		err := syscall.Mkfifo(pipeFile, 0666)
		if err != nil {
			log.Fatal("Make named pipe file error:", err)
		}
		f, err := os.OpenFile(pipeFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
		if err != nil {
			log.Fatalln(err)
		}
		for {
			select {
			case t := <-time.After(time.Second):
				today := fmt.Sprintf("%02d%02d%02d", t.Day(), t.Month(), t.Year()%100)
				ts := fmt.Sprintf("%02d%02d%02d.%03d", t.Hour(), t.Minute(), t.Second(), t.Nanosecond()/1000_000)
				nmeapayload := fmt.Sprintf("GPGGA,%s,4807.038,N,01131.000,E,1,08,0.9,545.4,M,46.9,M,,", ts)
				fullnmea := checksum(nmeapayload)
				_, err := f.WriteString(fullnmea)
				if err != nil {
					log.Fatalln(err)
				}
				fmt.Print(fullnmea)
				nmeapayload = fmt.Sprintf("GPRMC,%s.1,A,4807.038,N,01131.000,E,022.4,084.4,%s,,,A", ts, today)
				fullnmea = checksum(nmeapayload)
				_, err = f.WriteString(fullnmea)
				if err != nil {
					log.Fatalln(err)
				}
				fmt.Print(fullnmea)
			case <-s:
				return nil
			}
		}
	},
}
