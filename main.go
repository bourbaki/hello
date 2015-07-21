package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func main() {
	http.HandleFunc("/", hello)
	http.HandleFunc("/weather/", func(w http.ResponseWriter, r *http.Request) {
		city := strings.SplitN(r.URL.Path, "/", 3)[2]
		
		data, err := query(city)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(data)
	})
	http.ListenAndServe(":8080", nil)
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello!"))
}

type weatherData struct {
	Name string `json:"name"`
	Main struct {
		Kelvin float64 `json:"temp"`
	} `json:"main"`
}

func query(city string) (weatherData, error) {
	resp, err := http.Get("http://api.openweathermap.org/data/2.5/weather?q=" + city)
	if err != nil {
		return weatherData{}, err
	}
	
	defer resp.Body.Close() // execute when we return from main function
	
	var d weatherData
	
	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		// decoder takes io.Reader interface
		return weatherData{}, err
	}
	
	return d, nil
}

type weatherProvider interface {
	temperature(city string) (float64, error) // in Kelvin, naturally
}

type OpenWeatherMap struct{}

func (w openWeatherMap) temperature(city string) (float64, error) {
  resp, err := http.Get("http://api.openweathermap.org/data/2.5/weather?q=" + city)
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

  log.Printf("openWeatherMap: %s: %.2f", city, d.Main.Kelvin)
  return d.Main.Kelvin, nil
}