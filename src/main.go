package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "time"
)

// Struktura reprezentująca lokalizację z nazwą kraju, miasta i współrzędnymi geograficznymi

type Location struct {
    Country    string  `json:"country"`
    City       string  `json:"city"`
    Latitude   float64 `json:"latitude"`
    Longitude  float64 `json:"longitude"`
}

// Struktura dopasowana do formatu danych zwracanego przez API pogodowe

type WeatherResponse struct {
    CurrentWeather struct {
        Temperature float64 `json:"temperature"`
    } `json:"current_weather"`
}



var author = "Michał Krocz"
var port = "8080"

var locations = []Location{
	{"Polska", "Warszawa", 52.23, 21.01},
	{"Polska", "Lublin", 51.14, 22.34},
	{"Polska", "Kraków", 49.58, 19.47},
    {"Niemcy", "Berlin", 52.52, 13.41},
	{"Niemcy", "Hamburg", 53.55, 10.00},
	{"Niemcy", "Monachium", 48.08, 11.35},    
    {"USA", "Nowy Jork", 40.71, -74.01},
	{"USA", "Los Angeles", 34.03, -118.15},
	{"USA", "Chicago", 41.54, -87.39},
}

func main() {
	// Logowanie uruchomienia aplikacji
    now := time.Now().Format(time.RFC3339)
    log.Printf("Aplikacja uruchomiona o %s na porcie %s, autor %s", now, port, author)

    http.HandleFunc("/locations", handleLocations) // Zwraca dostępne lokalizacje
    http.HandleFunc("/weather", handleWeather) // Zwraca aktualną pogodę
	
	// Serwowanie plików statycznych z katalogu static
	http.Handle("/", http.FileServer(http.Dir("./static")))

	// Uruchomienie serwera HTTP
    log.Fatal(http.ListenAndServe(":"+port, nil))
}

// Zwracanie listy lokalizacji w formacie JSON
func handleLocations(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(locations)
}


// Pobieranie pogody dla wybranej lokalizacji
func handleWeather(w http.ResponseWriter, r *http.Request) {
    city := r.URL.Query().Get("city")
    country := r.URL.Query().Get("country")

    var selected *Location
    for _, loc := range locations {
        if loc.City == city && loc.Country == country {
            selected = &loc
            break
        }
    }

    if selected == nil {
        http.Error(w, "Location not found", http.StatusNotFound)
        return
    }

    // Składanie URL z parametrami
    url := fmt.Sprintf(
    "https://api.open-meteo.com/v1/forecast?latitude=%.2f&longitude=%.2f&current_weather=true",
    selected.Latitude, selected.Longitude,
)

	// Wysłanie zapytania do API pogodowego
    resp, err := http.Get(url)
    if err != nil {
        http.Error(w, "Failed to fetch weather data", http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()

	// Dekodowanie odpowiedzi JSON do struktury WeatherResponse
    var weather WeatherResponse
err = json.NewDecoder(resp.Body).Decode(&weather)
if err != nil {
    http.Error(w, "Failed to parse weather data", http.StatusInternalServerError)
    return
}

msg := fmt.Sprintf("Aktualna temperatura w %s, %s: %.1f°C",
    selected.City, selected.Country, weather.CurrentWeather.Temperature)

w.Header().Set("Content-Type", "text/plain")
w.Write([]byte(msg))

}
