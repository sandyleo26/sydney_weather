package main_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sandyleo26/sydney_weather/open_weather_map"
	"github.com/sandyleo26/sydney_weather/yahoo"
)

//
type mockYahooClient struct {
	Response *yahoo.Response
	Error    error
}

func (m mockYahooClient) Get() (*yahoo.Response, error) {
	return m.Response, m.Error
}

//
type mockOpenWeatherMapClient struct {
	Response *open_weather_map.Response
	Error    error
}

func (m mockOpenWeatherMapClient) Get() (*open_weather_map.Response, error) {
	return m.Response, m.Error
}

//
func f2c(f int) int {
	return int((f - 32) * 5 / 9)
}

func k2c(k float32) int {
	return int(k - 273.15)
}

func TestQueryYahooSuccess(t *testing.T) {
	mock := mockYahooClient{
		Response: &yahoo.Response{
			Query: yahoo.Query{
				Results: yahoo.Results{
					Channel: yahoo.Channel{
						Wind: yahoo.Wind{
							Speed: "99",
						},
						Item: yahoo.Item{
							Condition: yahoo.Condition{
								Temp: "70",
							},
						},
					},
				},
			},
		},
		Error: nil,
	}

	resp, err := yahoo.QueryYahoo(mock)
	assert.Nil(t, err)
	assert.Equal(t, 99, resp.WindSpeed)
	assert.Equal(t, f2c(70), resp.TemperatureDegrees)
}

func TestQueryOpenWeatherMapSuccess(t *testing.T) {
	mock := mockOpenWeatherMapClient{
		Response: &open_weather_map.Response{
			Wind: open_weather_map.Wind{
				Speed: 3.14,
			},
			Main: open_weather_map.Main{
				Temp: 300.88,
			},
		},
		Error: nil,
	}

	resp, err := open_weather_map.Query(mock)
	assert.Nil(t, err)
	assert.Equal(t, k2c(300.88), resp.TemperatureDegrees)
	assert.Equal(t, 3, resp.WindSpeed)
}
