package main

import (
	"flag"
	"fynecv/appdata"
	"log"
	"strconv"
	"time"
)

func main() {
	var (
		seconds = 60
		err     error
	)

	flag.Parse()
	if len(flag.Args()) > 0 {
		seconds, err = strconv.Atoi(flag.Args()[0])
		if err != nil {
			seconds = 60
		}
	}

	log.Println("Will run for", seconds, "seconds.")

	data := appdata.NewAppData()
	data.Monitor()

	timr := time.NewTimer(time.Second * time.Duration(seconds))
	<-timr.C

	data.StopMonitor()
}
