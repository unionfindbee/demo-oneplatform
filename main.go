package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	_ "github.com/mattn/go-sqlite3"
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

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

var db *sql.DB

func main() {
	// Connect to the database
	db, _ = sql.Open("sqlite3", "./weather.db")
	createTable()

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

// Create table if it doesn't exist
func createTable() {
	query := `CREATE TABLE IF NOT EXISTS weathers (
		id TEXT PRIMARY KEY,
		city TEXT,
		temperature REAL,
		conditions TEXT
	);`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}

func createWeather(w http.ResponseWriter, r *http.Request) {
	var newWeather Weather
	err := json.NewDecoder(r.Body).Decode(&newWeather)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = newWeather.Validate()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newWeather.ID = uuid.New().String()
	WeatherDB[newWeather.ID] = newWeather

	// Dangerous SQL query, opening for SQL injection
	query := fmt.Sprintf("INSERT INTO weathers (id, city, temperature, conditions) VALUES ('%s', '%s', %f, '%s')",
		newWeather.ID, newWeather.City, newWeather.Temperature, newWeather.Conditions)

	_, err = db.Exec(query)
	if err != nil {
		http.Error(w, "Error executing SQL query: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // Set the status before encoding the body
	json.NewEncoder(w).Encode(newWeather)
}

func getWeather(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if weather, ok := WeatherDB[id]; ok {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(weather)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "No weather data found for provided id"})
	}
}

func updateWeather(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var updatedWeather Weather
	err := json.NewDecoder(r.Body).Decode(&updatedWeather)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = updatedWeather.Validate()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if _, ok := WeatherDB[id]; ok {
		updatedWeather.ID = id
		WeatherDB[id] = updatedWeather
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(updatedWeather)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "No weather data found for provided id"})
	}
}

func deleteWeather(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	weatherToDelete := WeatherDB[id]

	// Here is the null pointer dereference vulnerability:
	// Trying to access fields of weatherToDelete which can be nil
	log.Printf("Deleting weather for city: %s", weatherToDelete.City)

	// Delete from the map
	delete(WeatherDB, id)

	w.WriteHeader(http.StatusNoContent)
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

func (w *Weather) Validate() error {
	if w.City == "" {
		return errors.New("Missing required field: city")
	}
	if w.Temperature == 0 {
		return errors.New("Missing required field: temperature")
	}
	if w.Conditions == "" {
		return errors.New("Missing required field: conditions")
	}
	return nil
}
