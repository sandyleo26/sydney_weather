package open_weather_map

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sandyleo26/sydney_weather/common"
)

//Response is a partial response returned from open weather map api
type Response struct {
	Base       string `json:"base"`
	Main       Main   `json:"main"`
	Visibility int    `json:"visibility"`
	Wind       Wind   `json:"wind"`
	DT         int    `json:"dt"`
	ID         int    `json:"id"`
	Name       string `json:"name"`
}

type Main struct {
	Temp float32 `json:"temp"`
}

type Wind struct {
	Speed float32 `json:"speed"`
}

//k2c converts Kelvin to Celsius
func k2c(k float32) int {
	return int(k - 273.15)
}

type Client interface {
	Get() (*Response, error)
}

type RealClient struct{}

//Get call open weather map api to retrieve weather info
func (c RealClient) Get() (*Response, error) {
	url := "http://api.openweathermap.org/data/2.5/weather?q=sydney,AU&appid=2326504fb9b100bee21400190e4dbe6d"

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Request to open weather map failed with error %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Request to open weather map failed with status %v", resp.StatusCode)
	}

	response := Response{}
	errDecode := json.NewDecoder(resp.Body).Decode(&response)
	if errDecode != nil {
		return nil, fmt.Errorf("Failed to decode response from open weather map with err %v", errDecode)
	}

	return &response, nil
}

//Query retrieve weather from open weather map
func Query(c Client) (*common.WeatherResponse, error) {
	resp, err := c.Get()
	if err != nil {
		return nil, err
	}

	return &common.WeatherResponse{
		WindSpeed:          int(resp.Wind.Speed),
		TemperatureDegrees: k2c(resp.Main.Temp),
	}, nil
}
