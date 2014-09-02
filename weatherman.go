package weatherman

import (
	"encoding/json"
	"log"
	"net/http"
)

type WeatherData struct {
	City string  `json:"city"`
	Temp float64 `json:"temp"`
	Took string  `json:"took"`
}

type WeatherProvider interface {
	Temperature(city string) (float64, error) //in Kelvin
}

type WeatherAggregator []WeatherProvider

func (w WeatherAggregator) Temperature(city string) (float64, error) {
	sum := 0.0
	for _, provider := range w {
		k, err := provider.Temperature(city)
		if err != nil {
			return 0, err
		}
		sum += k
	}

	return sum / float64(len(w)), nil
}

type OpenWeatherMap struct{}

func (w OpenWeatherMap) Temperature(city string) (float64, error) {
	endpoint := "http://api.openweathermap.org/data/2.5/weather?q="
	resp, err := http.Get(endpoint + city)

	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()

	var d struct {
		Main struct {
			Kelvin float64 `json:"temp"`
		} `json:"main"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return 0, err
	}

	log.Printf("openWeatherMap %s: %.2f", city, d.Main.Kelvin)
	return d.Main.Kelvin, nil
}

type WeatherUnderground struct {
	ApiKey string
}

func (w WeatherUnderground) Temperature(city string) (float64, error) {
	endpoint := "http://api.wunderground.com/api/" + w.ApiKey + "/conditions/q/"
	resp, err := http.Get(endpoint + city + ".json")
	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()

	var d struct {
		Observation struct {
			Celsius float64 `json:"temp_c"`
		} `json:"current_observation"`
	}

	kelvin := d.Observation.Celsius + 273.15
	log.Printf("weatherUnderground: %s: %.2f", city, kelvin)
	return kelvin, nil
}
