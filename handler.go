package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
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

	weatherResponse, err := QueryYahoo()
	fmt.Println("weatherResponse", weatherResponse)
	if err != nil {
		http.Error(w, "weather infomation is not available right now", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(weatherResponse)
}
