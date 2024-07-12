package main

import (
	"flag"
	"fynecv/appdata"
	"log"
	"strconv"
	"time"

	"fyne.io/fyne/v2/data/binding"
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
	data.Ready.AddListener(binding.NewDataListener(func() {
		isLoaded, _ := data.Ready.Get()
		if isLoaded {
			log.Println("DATA LOADED", isLoaded)
		}
	}))
	data.Monitor()

	<-time.NewTimer(time.Second * time.Duration(seconds)).C
	data.StopMonitor()
}
