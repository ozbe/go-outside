package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	lat, long, err := coordinates()
	if err != nil {
		log.Fatal(err)
	}

	err = weather(lat, long)
	if err != nil {
		log.Fatal(err)
	}
}

func weather(lat float64, long float64) error {
	apiKey := os.Getenv("OPENWEATHER_API_KEY")
	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?lat=%f&lon=%f&units=metric&appid=%s", lat, long, apiKey)

	res, err := http.Get(url)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	fmt.Println(string(body))

	return nil
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
