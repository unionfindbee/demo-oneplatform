package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	var err error
	db, err = sql.Open("sqlite3", "./weather_test.db")
	if err != nil {
		log.Fatal(err)
	}
	createTable()

	code := m.Run()

	db.Close()
	os.Exit(code)
}

func TestCreateWeather(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/weather", createWeather).Methods("POST")

	newWeather := &Weather{
		City:        "TestCity",
		Temperature: 20.5,
		Conditions:  "Cloudy",
	}

	jsonWeather, _ := json.Marshal(newWeather)
	request, _ := http.NewRequest("POST", "/weather", bytes.NewBuffer(jsonWeather))
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	assert.Equal(t, 201, response.Code, "Expected response code to be 201")
	var responseWeather Weather
	json.Unmarshal(response.Body.Bytes(), &responseWeather)
	assert.Equal(t, newWeather.City, responseWeather.City, "The city name should be the same")
}

func TestGetWeather(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/weather/{id}", getWeather).Methods("GET")

	// Create a new weather report
	newWeather := &Weather{
		ID:          "test-id",
		City:        "TestCity",
		Temperature: 20.5,
		Conditions:  "Cloudy",
	}
	WeatherDB[newWeather.ID] = *newWeather

	request, _ := http.NewRequest("GET", "/weather/test-id", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	assert.Equal(t, 200, response.Code, "Expected response code to be 200")
	var responseWeather Weather
	json.Unmarshal(response.Body.Bytes(), &responseWeather)
	assert.Equal(t, newWeather.ID, responseWeather.ID, "The weather ID should be the same")
}

func TestUpdateWeather(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/weather/{id}", updateWeather).Methods("PUT")

	// Create a new weather report
	newWeather := &Weather{
		ID:          "test-id",
		City:        "TestCity",
		Temperature: 20.5,
		Conditions:  "Cloudy",
	}
	WeatherDB[newWeather.ID] = *newWeather

	updatedWeather := &Weather{
		ID:          "test-id",
		City:        "UpdatedCity",
		Temperature: 22.5,
		Conditions:  "Sunny",
	}
	jsonWeather, _ := json.Marshal(updatedWeather)
	request, _ := http.NewRequest("PUT", "/weather/test-id", bytes.NewBuffer(jsonWeather))
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	assert.Equal(t, 200, response.Code, "Expected response code to be 200")
	var responseWeather Weather
	json.Unmarshal(response.Body.Bytes(), &responseWeather)
	assert.Equal(t, updatedWeather.City, responseWeather.City, "The city name should be updated")
}

func TestDeleteWeather(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/weather/{id}", deleteWeather).Methods("DELETE")

	// Create a new weather report
	newWeather := &Weather{
		ID:          "test-id",
		City:        "TestCity",
		Temperature: 20.5,
		Conditions:  "Cloudy",
	}
	WeatherDB[newWeather.ID] = *newWeather

	request, _ := http.NewRequest("DELETE", "/weather/test-id", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	assert.Equal(t, 204, response.Code, "Expected response code to be 204")
	_, exist := WeatherDB[newWeather.ID]
	assert.False(t, exist, "The weather report should be deleted")
}
