package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/sandyleo26/sydney_weather/open_weather_map"
	"github.com/sandyleo26/sydney_weather/yahoo"
)

//WeatherHandler http handler for weather endpoing
func WeatherHandler(w http.ResponseWriter, r *http.Request) {
	if !strings.EqualFold(r.Method, "GET") {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	queryValues := r.URL.Query()
	if !strings.EqualFold(queryValues.Get("city"), "sydney") {
		http.Error(w, "Sorry, only Sydney is supported at the moment", http.StatusBadRequest)
		return
	}

	resp, err := yahoo.QueryYahoo()
	if err != nil {
		log.Println(err.Error())
		log.Println("Fall back to open weather map...")
		resp, err = open_weather_map.Query()
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "weather infomation is not available right now", http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
