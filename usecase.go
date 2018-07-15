package main

import (
	"fmt"
	"log"

	"github.com/sandyleo26/sydney_weather/common"
	"github.com/sandyleo26/sydney_weather/open_weather_map"
	"github.com/sandyleo26/sydney_weather/yahoo"
)

// var cache struct {
// 	Time time.Time
// 	common.WeatherResponse
// }

// var mutex = &sync.Mutex{}

func GetWeather(yc yahoo.Client, oc open_weather_map.Client) (*common.WeatherResponse, error) {
	// if time.Now().Before(cache.Time.Add(time.Second * 3)) {
	// 	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	// 	w.WriteHeader(http.StatusOK)
	// 	json.NewEncoder(w).Encode(cache.WeatherResponse)
	// 	return
	// }

	resp, err := yahoo.QueryYahoo(yc)
	if err != nil {
		log.Println(err.Error())
		log.Println("Fall back to open weather map...")
		resp, err = open_weather_map.Query(oc)
		if err != nil {
			log.Println(err.Error())
			return nil, fmt.Errorf("weather infomation is not available right now")
		}
	}

	// mutex.Lock()
	// now := time.Now()
	// if now.After(cache.Time) {
	// 	cache.Time = now
	// 	cache.WeatherResponse = *resp
	// }
	// mutex.Unlock()

	return resp, nil
}
