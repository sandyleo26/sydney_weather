package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

type YahooResponseQueryResultsChannelItemCondition struct {
	Code string `json:"code"`
	Date string `json:"date"`
	Temp string `json:"temp"`
	Text string `json:"text"`
}

type YahooResponseQueryResultsChannelWind struct {
	Chill     string `json:"chill"`
	Direction string `json:"direction"`
	Speed     string `json:"speed"`
}

type YahooResponseQueryResultsChannelItem struct {
	Condition YahooResponseQueryResultsChannelItemCondition `json:"condition"`
}

type YahooResponseQueryResultsChannel struct {
	Wind YahooResponseQueryResultsChannelWind `json:"wind"`
	Item YahooResponseQueryResultsChannelItem `json:"item"`
}

type YahooResponseQueryResults struct {
	Channel YahooResponseQueryResultsChannel `json:"channel`
}

type YahooResponseQuery struct {
	Count   int                       `json:"count"`
	Created string                    `json:"created"`
	Lang    string                    `json:"lang"`
	Results YahooResponseQueryResults `json:"results"`
}

type YahooResponse struct {
	Query YahooResponseQuery `json:"query"`
}

//WeatherResponse response for weather handler
type WeatherResponse struct {
	WindSpeed          int `json:"wind_speed"`
	TemperatureDegrees int `json:"temperature_degrees"`
}

//WeatherHandler http handler for weather endpoing
func WeatherHandler(w http.ResponseWriter, r *http.Request) {
	if !strings.EqualFold(r.Method, "GET") {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
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

//QueryYahoo call yahoo weather API
func QueryYahoo() (*WeatherResponse, error) {
	yahooURL := "https://query.yahooapis.com/v1/public/yql?q=select%20item.condition%2C%20wind%20from%20weather.forecast%20where%20woeid%20%3D%201105779&format=json&env=store%3A%2F%2Fdatatables.org%2Falltableswithkeys"
	resp, err := http.Get(yahooURL)
	if err != nil {
		return nil, fmt.Errorf("Request to Yahoo failed with error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request to Yahoo failed with status: %v", resp.StatusCode)
	}

	yahooResponse := YahooResponse{}
	errDecode := json.NewDecoder(resp.Body).Decode(&yahooResponse)
	if errDecode != nil {
		return nil, fmt.Errorf("Failed to decode response from Yahoo with %v", errDecode)
	}

	windSpeed, windSpeedErr := strconv.Atoi(yahooResponse.Query.Results.Channel.Wind.Speed)
	if windSpeedErr != nil {
		return nil, fmt.Errorf("Failed to retrieve wind speed info with error %v", windSpeedErr)
	}

	temperatureDegrees, temperatureDegreesErr := strconv.Atoi(yahooResponse.Query.Results.Channel.Item.Condition.Temp)
	if temperatureDegreesErr != nil {
		return nil, fmt.Errorf("Failed to retrieve temperature info with error %v", temperatureDegreesErr)
	}

	return &WeatherResponse{
		WindSpeed:          windSpeed,
		TemperatureDegrees: temperatureDegrees,
	}, nil

}

func main() {

	router := mux.NewRouter()
	router.HandleFunc("/v1/weather", WeatherHandler)
	http.Handle("/", router)

	log.Println("Start listenning on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Println("Web server failed to start at 8080 witht error: ", err)
	}

}
