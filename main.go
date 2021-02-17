package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/common-nighthawk/go-figure"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	lat, long, err := coordinates()
	if err != nil {
		log.Fatal(err)
	}

	w, err := weather(lat, long)
	if err != nil {
		log.Fatal(err)
	}

	print(w)
	fmt.Scanln()
}

type GeoIpResponse struct {
	Latitude  float64 `json:"latitude`
	Longitude float64 `json:"longitude`
}

func coordinates() (lat float64, long float64, err error) {
	url := "https://freegeoip.app/json/"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	var result GeoIpResponse
	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return
	}

	lat = result.Latitude
	long = result.Longitude
	return
}

type WeatherResponse struct {
	Weather []Weather `json:"weather"`
	Main    Main      `json:"main"`
}

type Main struct {
	FeelsLike float64 `json:"feels_like"`
}

type MainWeatherCondition int

const (
	Thunderstorm MainWeatherCondition = iota
	Drizzle
	Rain
	Snow
	Mist
	Smoke
	Haze
	Dust
	Fog
	Sand
	Ash
	Squall
	Tornado
	Clear
	Clouds
)

func (m MainWeatherCondition) String() string {
	return [...]string{"Thunderstorm", "Drizzle", "Rain", "Snow", "Mist", "Smoke", "Haze", "Dust", "Fog", "Sand", "Ash", "Squall", "Tornado", "Clear", "Clouds"}[m]
}

func (m *MainWeatherCondition) UnmarshalJSON(b []byte) error {
	var s string
	json.Unmarshal(b, &s)
	switch s {
	case "Thunderstorm":
		*m = Thunderstorm
	case "Drizzle":
		*m = Drizzle
	case "Rain":
		*m = Rain
	case "Snow":
		*m = Snow
	case "Mist":
		*m = Mist
	case "Smoke":
		*m = Smoke
	case "Haze":
		*m = Haze
	case "Dust":
		*m = Dust
	case "Sand":
		*m = Sand
	case "Ash":
		*m = Ash
	case "Squall":
		*m = Squall
	case "Tornado":
		*m = Tornado
	case "Clear":
		*m = Clear
	case "Clouds":
		*m = Clouds
	default:
		return errors.New("Unknown main weather condition")
	}

	return nil
}

type Weather struct {
	Main MainWeatherCondition `json:"main"`
}

func weather(lat float64, long float64) (*WeatherResponse, error) {
	apiKey := os.Getenv("OPENWEATHER_API_KEY")
	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?lat=%f&lon=%f&units=metric&appid=%s", lat, long, apiKey)

	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	var result WeatherResponse
	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func print(w *WeatherResponse) {
	s := fmt.Sprintf("%.1f - %s", w.Main.FeelsLike, w.Weather[0].Main.String())
	f := figure.NewFigure(s, "", false)
	f.Print()
}
