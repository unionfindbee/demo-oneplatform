package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type Weather struct {
	ID          string  `json:"id"`
	City        string  `json:"city"`
	Temperature float64 `json:"temperature"`
	Conditions  string  `json:"conditions"`
}

var WeatherDB = make(map[string]Weather)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/weather", createWeather).Methods("POST")
	r.HandleFunc("/weather/{id}", getWeather).Methods("GET")
	r.HandleFunc("/weather/{id}", updateWeather).Methods("PUT")
	r.HandleFunc("/weather/{id}", deleteWeather).Methods("DELETE")
	r.HandleFunc("/weather-stream", weatherStream)
	http.ListenAndServe(":7070", r)
}

func createWeather(w http.ResponseWriter, r *http.Request) {
	var newWeather Weather
	json.NewDecoder(r.Body).Decode(&newWeather)

	// Here we generate a new unique ID
	newWeather.ID = uuid.New().String()
	WeatherDB[newWeather.ID] = newWeather

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newWeather)
}

func getWeather(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if weather, ok := WeatherDB[id]; ok {
		json.NewEncoder(w).Encode(weather)
		return
	}
	w.WriteHeader(http.StatusNotFound)
}

func updateWeather(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var updatedWeather Weather
	json.NewDecoder(r.Body).Decode(&updatedWeather)
	if _, ok := WeatherDB[id]; ok {
		updatedWeather.ID = id
		WeatherDB[id] = updatedWeather
		json.NewEncoder(w).Encode(updatedWeather)
		return
	}
	w.WriteHeader(http.StatusNotFound)
}

func deleteWeather(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if _, ok := WeatherDB[id]; ok {
		delete(WeatherDB, id)
		w.WriteHeader(http.StatusNoContent)
		return
	}
	w.WriteHeader(http.StatusNotFound)
}

func weatherStream(w http.ResponseWriter, r *http.Request) {
	conn, _ := upgrader.Upgrade(w, r, nil)
	for {
		// Assume we have some mechanism to get the latest weather update
		weatherUpdate := getLatestWeatherUpdate()
		if err := conn.WriteJSON(weatherUpdate); err != nil {
			return
		}
	}
}

func getLatestWeatherUpdate() Weather {
	// Return random weather from WeatherDB
	for _, weather := range WeatherDB {
		return weather
	}
	return Weather{}
}
