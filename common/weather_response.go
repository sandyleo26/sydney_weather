package common

//WeatherResponse response for weather handler
type WeatherResponse struct {
	WindSpeed          int `json:"wind_speed"`
	TemperatureDegrees int `json:"temperature_degrees"`
}
