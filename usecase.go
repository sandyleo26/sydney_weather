package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/sandyleo26/sydney_weather/common"
	"github.com/sandyleo26/sydney_weather/open_weather_map"
	"github.com/sandyleo26/sydney_weather/yahoo"
)

var cache struct {
	Time time.Time
	common.WeatherResponse
}

var mutex = &sync.Mutex{}

//GetWeather fetch weather info
func GetWeather(yc yahoo.Client, oc open_weather_map.Client, useCache bool) (*common.WeatherResponse, error) {
	if useCache && time.Now().Before(cache.Time.Add(time.Second*3)) {
		return &cache.WeatherResponse, nil
	}

	resp, err := yahoo.QueryYahoo(yc)
	if err != nil {
		log.Println(err.Error())
		log.Println("Fall back to open weather map...")
		resp, err = open_weather_map.Query(oc)
		if err != nil {
			log.Println(err.Error())
			if useCache {
				return &cache.WeatherResponse, nil
			}
			return nil, fmt.Errorf("weather infomation is not available right now")
		}
	}

	if useCache {
		mutex.Lock()
		now := time.Now()
		if now.After(cache.Time) {
			cache.Time = now
			cache.WeatherResponse = *resp
		}
		mutex.Unlock()
	}

	return resp, nil
}
