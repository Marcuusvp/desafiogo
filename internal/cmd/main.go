package main

import (
	"desafiofc/internal/domain"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	eventsMutex sync.RWMutex
	events      []domain.Event
	spotsMutex  sync.RWMutex
	spots       []domain.Spot
)

func main() {
	loadData()

	http.HandleFunc("/events", listEventsHandler)
	http.HandleFunc("/events/", eventHandler)
	http.HandleFunc("/event/", reserveSpotHandler)

	fmt.Println("Servidor iniciado em http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func loadData() {
	jsonData, err := ioutil.ReadFile("db.json")
	if err != nil {
		log.Fatalf("Erro ao ler arquivo db.json: %v", err)
	}

	var data map[string]interface{}
	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		log.Fatalf("Erro ao decodificar JSON: %v", err)
	}

	eventsData, ok := data["events"].([]interface{})
	if !ok {
		log.Fatal("Dados de eventos não encontrados ou no formato incorreto")
	}
	for _, eventData := range eventsData {
		eventMap := eventData.(map[string]interface{})
		event := domain.Event{
			ID:          int(eventMap["id"].(float64)),
			Name:        eventMap["name"].(string),
			Organization: eventMap["organization"].(string),
			Date:        parseDate(eventMap["date"].(string)),
			Price:       int(eventMap["price"].(float64)),
			Rating:      eventMap["rating"].(string),
			ImageURL:    eventMap["image_url"].(string),
			CreatedAt:   parseDate(eventMap["created_at"].(string)),
			Location:    eventMap["location"].(string),
		}
		events = append(events, event)
	}

	spotsData, ok := data["spots"].([]interface{})
	if !ok {
		log.Fatal("Dados de spots não encontrados ou no formato incorreto")
	}
	for _, spotData := range spotsData {
		spotMap := spotData.(map[string]interface{})
		spot := domain.Spot{
			ID:      int(spotMap["id"].(float64)),
			Name:    spotMap["name"].(string),
			Status:  spotMap["status"].(string),
			EventID: int(spotMap["event_id"].(float64)),
		}
		spots = append(spots, spot)
	}
}

func parseDate(dateStr string) time.Time {
	layout := "2006-01-02T15:04:05"
	t, err := time.Parse(layout, dateStr)
	if err != nil {
		log.Fatalf("Erro ao converter data: %v", err)
	}
	return t
}

func listEventsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	eventsMutex.RLock()
	defer eventsMutex.RUnlock()

	jsonData, err := json.Marshal(events)
	if err != nil {
		http.Error(w, "Erro ao serializar eventos", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func eventHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/events/")
	parts := strings.Split(path, "/")

	if len(parts) == 1 {
		getEventHandler(w, r, parts[0])
	} else if len(parts) == 2 && parts[1] == "spots" {
		listEventSpotsHandler(w, r, parts[0])
	} else {
		http.NotFound(w, r)
	}
}

func getEventHandler(w http.ResponseWriter, r *http.Request, eventIDStr string) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	eventID, err := strconv.Atoi(eventIDStr)
	if err != nil {
		http.Error(w, "ID do evento inválido", http.StatusBadRequest)
		return
	}

	eventsMutex.RLock()
	defer eventsMutex.RUnlock()

	var event domain.Event
	for _, e := range events {
		if e.ID == eventID {
			event = e
			break
		}
	}

	if event.ID == 0 {
		http.Error(w, "Evento não encontrado", http.StatusNotFound)
		return
	}

	jsonData, err := json.Marshal(event)
	if err != nil {
		http.Error(w, "Erro ao serializar evento", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func listEventSpotsHandler(w http.ResponseWriter, r *http.Request, eventIDStr string) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	eventID, err := strconv.Atoi(eventIDStr)
	if err != nil {
		http.Error(w, "ID do evento inválido", http.StatusBadRequest)
		return
	}

	spotsMutex.RLock()
	defer spotsMutex.RUnlock()

	var eventSpots []domain.Spot
	for _, s := range spots {
		if s.EventID == eventID {
			eventSpots = append(eventSpots, s)
		}
	}

	jsonData, err := json.Marshal(eventSpots)
	if err != nil {
		http.Error(w, "Erro ao serializar spots", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func reserveSpotHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/event/")
	parts := strings.Split(path, "/")

	if len(parts) != 2 || parts[1] != "reserve" {
		http.NotFound(w, r)
		return
	}

	eventIDStr := parts[0]
	eventID, err := strconv.Atoi(eventIDStr)
	if err != nil {
		http.Error(w, "ID do evento inválido", http.StatusBadRequest)
		return
	}

	var reqBody map[string][]string
	err = json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, "Erro ao ler body da requisição", http.StatusBadRequest)
		return
	}

	spotNames, ok := reqBody["spots"]
	if !ok || len(spotNames) == 0 {
		http.Error(w, "Formato de requisição inválido", http.StatusBadRequest)
		return
	}

	spotsMutex.Lock()
	defer spotsMutex.Unlock()

	seen := make(map[string]bool)
	for _, spotName := range spotNames {
		if seen[spotName] {
			http.Error(w, fmt.Sprintf("Você está tentando reservar mais de uma vez o Spot %s", spotName), http.StatusBadRequest)
			return
		}
		seen[spotName] = true
	}

	for _, spotName := range spotNames {
		found := false
		for _, spot := range spots {
			if spot.EventID == eventID && spot.Name == spotName {
				found = true
				if spot.Status != "available" {
					http.Error(w, fmt.Sprintf("Spot %s já reservado", spotName), http.StatusBadRequest)
					return
				}
			}
		}
		if !found {
			http.Error(w, fmt.Sprintf("Spot %s não encontrado para o evento %d", spotName, eventID), http.StatusNotFound)
			return
		}
	}

	for _, spotName := range spotNames {
		for i, spot := range spots {
			if spot.EventID == eventID && spot.Name == spotName {
				spots[i].Status = "reserved"
				break
			}
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

