package web

import (
	"fmt"
	"fynecv/vision"
	"log"
	"net/http"
)

func Serve(cameras []*vision.Camera) {
	const url = "192.168.0.7:9000"

	for i, camera := range cameras {
		path := "/"
		if i > 0 {
			path = fmt.Sprintf("/%d/", i)
		}
		http.Handle(path, camera.StreamHook.Stream)
		go camera.Serve()
		log.Printf("Serving %s @%s%s\n", camera.Name, url, path)
	}

	server := &http.Server{
		Addr:         url,
		ReadTimeout:  0,
		WriteTimeout: 0,
	}

	log.Fatal(server.ListenAndServe())

}
