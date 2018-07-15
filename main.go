package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	router := mux.NewRouter()
	router.HandleFunc("/v1/weather", WeatherHandler)
	http.Handle("/", router)

	log.Println("Start listenning on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Println("Web server failed to start at 8080 witht error: ", err)
	}

}
