package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/evandigby/rtb"
	"github.com/evandigby/rtb/amqp"
	"os"
)

func main() {
	var fileName string
	var ampqAddress string
	flag.StringVar(&ampqAddress, "amqp-address", "amqp://guest:guest@localhost:5672/", "AMQP Logging address.")
	flag.StringVar(&fileName, "filename", "", "File to write to (stdout if not specified)")
	flag.Parse()

	if fileName == "" {
		fmt.Println("No file name specified, logging to screen.")
	}
	var logConsumer rtb.BidLogConsumer

	logConsumer = amqp.NewAmqpBidLogger(ampqAddress, "rtb-interview")

	logChan := logConsumer.LogChannel()

	for logItem := range logChan {
		js, err := json.Marshal(logItem)

		if err != nil {
			return
		}
		if fileName == "" {
			fmt.Println(string(js[:]))
		} else {
			file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

			if err != nil {
				fmt.Println("Error writing to log file.")
			} else {
				fmt.Fprintln(file, string(js[:]), "\n")
				file.Close()
			}
		}
	}
}
