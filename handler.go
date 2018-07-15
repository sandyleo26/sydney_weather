package main

import (
	"encoding/json"
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

	resp, err := GetWeather(yahoo.RealClient{}, open_weather_map.RealClient{}, true)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
