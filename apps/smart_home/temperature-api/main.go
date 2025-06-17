// @title SmartHome Temperature API
// @version 1.0
// @description API для работы с датчиками температуры
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8081
// @BasePath /

package main

import (
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	httpSwagger "github.com/swaggo/http-swagger"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	_ "temperature-api/docs"
	"time"
)

type Sensor struct {
	ID       string  `json:"id" example:"sensor-001"`
	Name     string  `json:"name" example:"Living Room Sensor"`
	Type     string  `json:"type" example:"temperature"`
	Location string  `json:"location" example:"Living Room"`
	Value    float64 `json:"value" example:"22.5"`
	Unit     string  `json:"unit" example:"°C"`
	Status   string  `json:"status" example:"active"`
}

var sensorsDB = make(map[string]Sensor)
var producer sarama.SyncProducer

// @title SmartHome Temperature API
// @version 1.0
// @description API для работы с датчиками температуры
// @host localhost:8081
// @BasePath /
func main() {
	// Сначала регистрируем специфичные пути
	http.HandleFunc("/health", healthCheck)
	http.HandleFunc("/temperature", temperatureHandler)
	http.HandleFunc("/api/v1/sensors", handleSensors)
	http.HandleFunc("/api/v1/sensors/", handleSensorRoutes)
	http.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			w.Write([]byte("Temperature API is running. Use /temperature endpoint"))
			return
		}
		http.NotFound(w, r)
	})

	http.ListenAndServe(":8081", nil)

	brokers := strings.Split(os.Getenv("KAFKA_BROKERS"), ",")
	if err := initKafka(brokers); err != nil {
		log.Fatalf("Failed to initialize Kafka producer: %v", err)
	}
	defer producer.Close()
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
// handleSensors godoc
// @Summary Работа с сенсорами
// @Description Создание сенсора или получение списка сенсоров
// @Tags Сенсоры
// @Accept json
// @Produce json
// @Success 200 {array} Sensor "Список сенсоров"
// @Success 201 {object} Sensor "Созданный сенсор"
// @Failure 400 {object} map[string]string "Неверный запрос"
// @Failure 405 {object} map[string]string "Метод не разрешен"
// @Router /api/v1/sensors [get]
// @Router /api/v1/sensors [post]
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

// @Summary Получить текущую температуру
// @Description Возвращает текущую температуру для указанного местоположения или датчика
// @Tags Телеметрия
// @Accept  json
// @Produce  json
// @Param location query string false "Название местоположения"
// @Param sensorId query string false "Идентификатор датчика"
// @Success 200 {object} TemperatureResponse
// @Failure 400 {object} map[string]string
// @Router /temperature [get]
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

	publishTemperatureEvent(sensorID, temp)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "OK"})
}

// handleSensorByID godoc
// @Summary Работа с конкретным сенсором
// @Description Получение, обновление или удаление сенсора по ID
// @Tags Сенсоры
// @Accept json
// @Produce json
// @Param id path string true "ID сенсора"
// @Success 200 {object} Sensor "Информация о сенсоре"
// @Success 204 "Сенсор удален"
// @Failure 400 {object} map[string]string "Неверный запрос"
// @Failure 404 {object} map[string]string "Сенсор не найден"
// @Failure 405 {object} map[string]string "Метод не разрешен"
// @Router /api/v1/sensors/{id} [get]
// @Router /api/v1/sensors/{id} [put]
// @Router /api/v1/sensors/{id} [delete]
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

func initKafka(brokers []string) error {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	var err error
	producer, err = sarama.NewSyncProducer(brokers, config)
	return err
}

func publishTemperatureEvent(sensorID string, temp float64) error {
	topic := "temperature_events"
	event := map[string]interface{}{
		"sensorId":    sensorID,
		"temperature": temp,
		"timestamp":   time.Now().UTC(),
	}
	jsonData, _ := json.Marshal(event)

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(jsonData),
	}

	_, _, err := producer.SendMessage(msg)
	return err
}
