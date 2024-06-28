package web

import (
	"fmt"
	"fynecv/cv"
	"log"
	"net/http"
)

func Serve(devices []*cv.Device) {
	// var (
	// 	webCam = cv.NewDevice(0, gocv.VideoCaptureV4L)
	// 	ipCam  = cv.NewDevice("http://192.168.0.25:8080", gocv.VideoCaptureAny)
	// )

	for i, device := range devices {
		path := "/"
		if i > 0 {
			path = fmt.Sprintf("/%d/", i)
		}
		http.Handle(path, device.StreamHook.Stream)
		go device.Capture()
	}

	url := "192.168.0.7:9000"
	fmt.Println("Capturing. Point your browser to " + url)

	server := &http.Server{
		Addr:         url,
		ReadTimeout:  0,
		WriteTimeout: 0,
	}

	log.Fatal(server.ListenAndServe())

}
