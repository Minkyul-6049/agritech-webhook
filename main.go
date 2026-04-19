package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os" // Add environment variable reader

	"github.com/joho/godotenv" // Add library to read .env files
)

// Change to global variables (use var instead of const)
var TelegramBotToken string
var TelegramChatID string

// 1. Grafana Alert Payload structure
type GrafanaAlert struct {
	Status string `json:"status"` // firing or resolved
	Alerts []struct {
		Labels      map[string]string `json:"labels"`
		Annotations map[string]string `json:"annotations"`
	} `json:"alerts"`
}

func sendTelegramMessage(text string) {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", TelegramBotToken)
	payload := map[string]string{"chat_id": TelegramChatID, "text": text}
	jsonPayload, _ := json.Marshal(payload)
	http.Post(url, "application/json", bytes.NewBuffer(jsonPayload))
}

// 2. Webhook function for request treatment
func webhookHandler(w http.ResponseWriter, r *http.Request) {
	// POST request
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var alert GrafanaAlert
	// income json data grafana mapping
	err := json.NewDecoder(r.Body).Decode(&alert)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// 3. auto action balancing up to alerts
	if alert.Status == "firing" {
		log.Println("[Warning] Grafana Alert receive: farm temperature > 30C!")

		// Extract specific info (e.g., alertname) from the first alert in Grafana JSON
		alertName := alert.Alerts[0].Labels["alertname"]

		// Generate dynamic message
		msg := fmt.Sprintf("[ALERT] %s Triggered!\nFarm temperature exceeded 30C.\nAction: Cooling system operating...", alertName)

		sendTelegramMessage(msg)

	} else if alert.Status == "resolved" {
		log.Println("[Normal] Farm Temp Normal checked.")

		msg := "[Normal] Farm Temp Normal.\nAction: Cooling system stopped."
		sendTelegramMessage(msg)
	}
	// Grafana receiving 200 Ok
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Webhook received successfully")
}

func main() {
	// Load .env file first when server starts
	err := godotenv.Load()
	if err != nil {
		log.Println("[Warning] .env file not found. Using system environment variables.")
	}

	// Load values from environment variables
	TelegramBotToken = os.Getenv("TELEGRAM_BOT_TOKEN")
	TelegramChatID = os.Getenv("TELEGRAM_CHAT_ID")

	// Safety check: fail if token is missing
	if TelegramBotToken == "" || TelegramChatID == "" {
		log.Fatal("[Error] TELEGRAM_BOT_TOKEN or TELEGRAM_CHAT_ID is missing in .env file!")
	}

	// "/webhook" data income and webhookHandler function operating connection
	http.HandleFunc("/webhook", webhookHandler)

	port := "8080"
	log.Printf("[Info] Agritech Webhook Server is working on port %s...\n", port)

	// Server start (If Error, program off)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
