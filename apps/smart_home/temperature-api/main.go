package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

type Sensor struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Type     string  `json:"type"`
	Location string  `json:"location"`
	Value    float64 `json:"value"`
	Unit     string  `json:"unit"`
	Status   string  `json:"status"`
}

var sensorsDB = make(map[string]Sensor)

func main() {
	// Сначала регистрируем специфичные пути
	http.HandleFunc("/health", healthCheck)
	http.HandleFunc("/temperature", temperatureHandler)
	http.HandleFunc("/api/v1/sensors", handleSensors)
	http.HandleFunc("/api/v1/sensors/", handleSensorRoutes)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			w.Write([]byte("Temperature API is running. Use /temperature endpoint"))
			return
		}
		http.NotFound(w, r)
	})

	http.ListenAndServe(":8081", nil)
}

func handleSensorRoutes(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	basePath := "/api/v1/sensors/"
	id := strings.TrimPrefix(path, basePath)

	// Разделяем путь на компоненты
	parts := strings.Split(id, "/")

	// Если после /sensors/ ничего нет
	if len(parts) == 0 || parts[0] == "" {
		http.NotFound(w, r)
		return
	}

	sensorID := parts[0]

	// Если запрос заканчивается на /value
	if len(parts) > 1 && parts[1] == "value" {
		handleSensorValue(w, r, sensorID) // Передаем sensorID
		return
	}

	// Все остальные запросы к /api/v1/sensors/{id}
	handleSensorByID(w, r, sensorID) // Передаем sensorID
}

// Пример обработчика для /api/v1/sensors
func handleSensors(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		// Логика получения всех сенсоров
		sensors := make([]Sensor, 0, len(sensorsDB))
		for _, s := range sensorsDB {
			sensors = append(sensors, s)
		}
		json.NewEncoder(w).Encode(sensors)
	case "POST":
		// Логика создания сенсора
		var sensor Sensor
		if err := json.NewDecoder(r.Body).Decode(&sensor); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		sensor.ID = generateID()
		// Генерируем начальное значение температуры
		sensor.Value = 15.0 + rand.Float64()*15.0
		sensor.Status = "active"

		sensorsDB[sensor.ID] = sensor
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(sensor)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

type TemperatureResponse struct {
	Temperature float64 `json:"temperature"`
	Location    string  `json:"location"`
	SensorID    string  `json:"sensorId"`
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func temperatureHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	location := query.Get("location")
	sensorID := query.Get("sensorId")

	// Если не указаны ни location, ни sensorId - возвращаем ошибку
	if location == "" && sensorID == "" {
		http.Error(w, `{"error": "Please provide either location or sensorId"}`, http.StatusBadRequest)
		return
	}

	// Если location не указан, определяем его по sensorId
	if location == "" {
		switch sensorID {
		case "1":
			location = "Living Room"
		case "2":
			location = "Bedroom"
		case "3":
			location = "Kitchen"
		default:
			location = "Unknown"
		}
	}

	// Если sensorId не указан, определяем его по location
	if sensorID == "" {
		switch location {
		case "Living Room":
			sensorID = "1"
		case "Bedroom":
			sensorID = "2"
		case "Kitchen":
			sensorID = "3"
		default:
			sensorID = "0"
		}
	}

	// Генерируем случайную температуру от 15.0 до 30.0
	temp := 15.0 + rand.Float64()*15.0

	response := TemperatureResponse{
		Temperature: temp,
		Location:    location,
		SensorID:    sensorID,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "OK"})
}

func handleSensorByID(w http.ResponseWriter, r *http.Request, id string) {
	switch r.Method {
	case "GET":
		sensor, exists := sensorsDB[id]
		if !exists {
			http.NotFound(w, r)
			return
		}
		json.NewEncoder(w).Encode(sensor)

	case "PUT":
		var updatedSensor Sensor
		if err := json.NewDecoder(r.Body).Decode(&updatedSensor); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if _, exists := sensorsDB[id]; !exists {
			http.NotFound(w, r)
			return
		}

		updatedSensor.ID = id
		sensorsDB[id] = updatedSensor
		json.NewEncoder(w).Encode(updatedSensor)

	case "DELETE":
		if _, exists := sensorsDB[id]; !exists {
			http.NotFound(w, r)
			return
		}
		delete(sensorsDB, id)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Sensor deleted successfully",
			"id":      id,
		})

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func generateID() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%d", rand.Intn(1000)+1)
}

func handleSensorValue(w http.ResponseWriter, r *http.Request, id string) {
	if r.Method != "PATCH" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sensor, exists := sensorsDB[id]
	if !exists {
		http.NotFound(w, r)
		return
	}

	// Обновляем только value и status
	var updateData struct {
		Value  float64 `json:"value"`
		Status string  `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Обновляем только указанные поля
	sensor.Value = updateData.Value
	sensor.Status = updateData.Status
	sensorsDB[id] = sensor

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sensor)
}
