package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
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

//f2c converts Fahrenheit to Celsius
func f2c(f int) int {
	return (f - 32) * 5 / 9
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

	temp, tempErr := strconv.Atoi(yahooResponse.Query.Results.Channel.Item.Condition.Temp)
	if tempErr != nil {
		return nil, fmt.Errorf("Failed to retrieve temperature info with error %v", tempErr)
	}

	return &WeatherResponse{
		WindSpeed:          windSpeed,
		TemperatureDegrees: f2c(temp),
	}, nil

}
