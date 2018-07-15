package main

import (
	"fmt"
	"strconv"
	"testing"
	"time"

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

func newMockYahooClient() *mockYahooClient {
	return &mockYahooClient{
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
}

//
type mockOpenWeatherMapClient struct {
	Response *open_weather_map.Response
	Error    error
}

func (m mockOpenWeatherMapClient) Get() (*open_weather_map.Response, error) {
	return m.Response, m.Error
}

func newMockOpenWeatherMapClient() *mockOpenWeatherMapClient {
	return &mockOpenWeatherMapClient{
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
}

//
func f2c(f int) int {
	return int((f - 32) * 5 / 9)
}

func k2c(k float32) int {
	return int(k - 273.15)
}

func TestQueryYahooSuccess(t *testing.T) {
	m := *newMockYahooClient()

	resp, err := yahoo.QueryYahoo(m)
	assert.Nil(t, err)
	assert.Equal(t, m.Response.Query.Results.Channel.Wind.Speed, strconv.Itoa(resp.WindSpeed))
	temp, _ := strconv.Atoi(m.Response.Query.Results.Channel.Item.Condition.Temp)
	assert.Equal(t, f2c(temp), resp.TemperatureDegrees)
}

func TestQueryOpenWeatherMapSuccess(t *testing.T) {
	mock := *newMockOpenWeatherMapClient()

	resp, err := open_weather_map.Query(mock)
	assert.Nil(t, err)
	assert.Equal(t, k2c(mock.Response.Main.Temp), resp.TemperatureDegrees)
	assert.Equal(t, int(mock.Response.Wind.Speed), resp.WindSpeed)
}

func TestGetWeatherSuccess(t *testing.T) {
	my := *newMockYahooClient()
	mo := *newMockOpenWeatherMapClient()

	resp, err := GetWeather(my, mo, false)
	assert.Nil(t, err)
	assert.Equal(t, my.Response.Query.Results.Channel.Wind.Speed, strconv.Itoa(resp.WindSpeed))
	temp, _ := strconv.Atoi(my.Response.Query.Results.Channel.Item.Condition.Temp)
	assert.Equal(t, f2c(temp), resp.TemperatureDegrees)
}

func TestGetWeatherFailOver(t *testing.T) {
	my := *newMockYahooClient()
	my.Error = fmt.Errorf("yahoo is hacked")
	mo := *newMockOpenWeatherMapClient()

	resp, err := GetWeather(my, mo, false)
	assert.Nil(t, err)
	assert.Equal(t, k2c(mo.Response.Main.Temp), resp.TemperatureDegrees)
	assert.Equal(t, int(mo.Response.Wind.Speed), resp.WindSpeed)
}

func TestGetWeatherUseCache(t *testing.T) {
	mOld := *newMockYahooClient()

	resp, _ := GetWeather(mOld, mockOpenWeatherMapClient{}, true)
	assert.Equal(t, mOld.Response.Query.Results.Channel.Wind.Speed, strconv.Itoa(resp.WindSpeed))
	temp, _ := strconv.Atoi(mOld.Response.Query.Results.Channel.Item.Condition.Temp)
	assert.Equal(t, f2c(temp), resp.TemperatureDegrees)

	mNew := *newMockYahooClient()
	mNew.Response.Query.Results.Channel.Wind.Speed = "199"
	resp, _ = GetWeather(mNew, mockOpenWeatherMapClient{}, true)
	// less than 3 seconds so should use old
	assert.Equal(t, mOld.Response.Query.Results.Channel.Wind.Speed, strconv.Itoa(resp.WindSpeed))

	// after 3 seconds so should use new
	time.Sleep(4 * time.Second)
	resp, _ = GetWeather(mNew, mockOpenWeatherMapClient{}, true)
	assert.Equal(t, mNew.Response.Query.Results.Channel.Wind.Speed, strconv.Itoa(resp.WindSpeed))

	// both services are down; should use cache
	time.Sleep(4 * time.Second)
	myBad := *newMockYahooClient()
	myBad.Error = fmt.Errorf("yahoo is hacked")
	moBad := *newMockOpenWeatherMapClient()
	moBad.Error = fmt.Errorf("open wheather map is hacked")
	resp, _ = GetWeather(myBad, moBad, true)
	assert.Equal(t, mNew.Response.Query.Results.Channel.Wind.Speed, strconv.Itoa(resp.WindSpeed))
}
