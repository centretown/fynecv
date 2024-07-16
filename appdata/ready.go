package appdata

import (
	"log"
)

func (data *AppData) GetReady() {
	err := data.Monitor()
	if err != nil {
		log.Fatal(err)
	}

	// wait := make(chan int)
	// data.Ready.AddListener(binding.NewDataListener(func() {
	// 	if ready, _ := data.Ready.Get(); ready {
	// 		log.Println("STATE LOADED")
	// 	}
	// 	if data.Err != nil {
	// 		log.Fatal(data.Err)
	// 	}
	// 	wait <- 1
	// }))
	// <-wait
}
