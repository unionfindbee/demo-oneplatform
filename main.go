package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"os/signal"
	"syscall"

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
	// Create a context that we can cancel
	ctx, cancel := context.WithCancel(context.Background())

	r := mux.NewRouter()
	r.HandleFunc("/weather", createWeather).Methods("POST")
	r.HandleFunc("/weather/{id}", getWeather).Methods("GET")
	r.HandleFunc("/weather/{id}", updateWeather).Methods("PUT")
	r.HandleFunc("/weather/{id}", deleteWeather).Methods("DELETE")
	r.HandleFunc("/weather-stream", weatherStream)

	server := http.Server{
		Addr:    ":7070",
		Handler: r,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			// We expect errors to happen when the server is Shutdown or closed,
			// but not in other cases.
			if err != http.ErrServerClosed {
				log.Fatalf("ListenAndServe(): %s", err)
			}
		}
	}()

	// Setup the shutdown signal handler
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	// We received an interrupt/kill signal; shut down gracefully.
	log.Println("Gracefully shutting down...")
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Could not gracefully shutdown the server: %s", err)
	}
	cancel()

	log.Println("Server stopped")
}

func createWeather(w http.ResponseWriter, r *http.Request) {
	var newWeather Weather
	json.NewDecoder(r.Body).Decode(&newWeather)

	// Here we generate a new unique ID
	newWeather.ID = uuid.New().String()
	WeatherDB[newWeather.ID] = newWeather

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // Set the status before encoding the body
	json.NewEncoder(w).Encode(newWeather)
}

func getWeather(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if weather, ok := WeatherDB[id]; ok {
		json.NewEncoder(w).Encode(weather)
		return
	}
	w.Header().Set("Content-Type", "application/json")
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
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
}

func deleteWeather(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if _, ok := WeatherDB[id]; ok {
		delete(WeatherDB, id)
		w.WriteHeader(http.StatusNoContent)
		return
	}
	w.Header().Set("Content-Type", "application/json")
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

func Fuzz(data []byte) int {
	if len(data) < 1 {
		// Not enough data for a test
		return -1
	}

	// Choose an endpoint and method based on the first byte
	var routeName string
	var routeFunc func(http.ResponseWriter, *http.Request)
	switch data[0] % 5 {
	case 0:
		routeName = "/weather"
		routeFunc = createWeather
	case 1:
		routeName = "/weather/1"
		routeFunc = getWeather
	case 2:
		routeName = "/weather/1"
		routeFunc = updateWeather
	case 3:
		routeName = "/weather/1"
		routeFunc = deleteWeather
	case 4:
		routeName = "/weather-stream"
		routeFunc = weatherStream
	}

	// Create an HTTP request and response recorder
	r, err := http.NewRequest("GET", routeName, bytes.NewReader(data[1:]))
	if err != nil {
		return 0
	}
	w := httptest.NewRecorder()

	// Call the chosen function and check the response
	routeFunc(w, r)
	if w.Code == http.StatusOK || w.Code == http.StatusCreated {
		return 1
	}
	return 0
}
