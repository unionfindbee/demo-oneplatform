package main

import (
	"encoding/json"
	"net/http"
	"sync"
)

type Weather struct {
	ID          string  `json:"id,omitempty"`
	City        string  `json:"city"`
	Temperature float64 `json:"temperature"`
	Conditions  string  `json:"conditions"`
}

var (
	mu       sync.RWMutex
	weathers = make(map[string]Weather)
)

func main() {
	http.HandleFunc("/weather", weatherHandler)
	http.ListenAndServe(":7070", nil)
}

func weatherHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		var weather Weather
		if err := json.NewDecoder(r.Body).Decode(&weather); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		mu.Lock()
		weathers[weather.ID] = weather
		mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(weather)

	case "GET":
		id := r.URL.Path[len("/weather/"):]
		mu.RLock()
		weather, ok := weathers[id]
		mu.RUnlock()

		if !ok {
			http.Error(w, "Weather not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(weather)

	case "PUT":
		id := r.URL.Path[len("/weather/"):]
		var weather Weather
		if err := json.NewDecoder(r.Body).Decode(&weather); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		mu.Lock()
		_, ok := weathers[id]
		if !ok {
			mu.Unlock()
			http.Error(w, "Weather not found", http.StatusNotFound)
			return
		}

		weathers[id] = weather
		mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(weather)

	case "DELETE":
		id := r.URL.Path[len("/weather/"):]
		mu.Lock()
		_, ok := weathers[id]
		if !ok {
			mu.Unlock()
			http.Error(w, "Weather not found", http.StatusNotFound)
			return
		}

		delete(weathers, id)
		mu.Unlock()

		w.WriteHeader(http.StatusNoContent)

	default:
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
	}
}
