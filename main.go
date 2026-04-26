package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

// GrafanaAlert represents the expected JSON payload from Grafana webhooks
type GrafanaAlert struct {
	Title    string `json:"title"`
	RuleName string `json:"ruleName"`
	State    string `json:"state"`
	Message  string `json:"message"`
}

func main() {
	// 1. Load environment variables for Telegram credentials
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	chatID := os.Getenv("TELEGRAM_CHAT_ID")

	if botToken == "" || chatID == "" {
		log.Fatal("Error: TELEGRAM_BOT_TOKEN and TELEGRAM_CHAT_ID must be set in the environment.")
	}

	// 2. Define HTTP routes
	// Health check endpoint for Docker/Kubernetes liveness probes
	http.HandleFunc("/health", healthCheckHandler)
	
	// Main alert receiver endpoint (Make sure Grafana URL ends with /alert)
	http.HandleFunc("/alert", func(w http.ResponseWriter, r *http.Request) {
		handleAlert(w, r, botToken, chatID)
	})

	// 3. Start the HTTP server
	port := ":8080"
	log.Printf("Starting Agritech webhook server on port %s", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// healthCheckHandler responds with a simple 200 OK status
func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK: Agritech Webhook Server is running securely at the edge.\n"))
}

// handleAlert processes incoming Grafana alerts and forwards them to Telegram
func handleAlert(w http.ResponseWriter, r *http.Request, botToken string, chatID string) {
	// Ensure only POST requests are accepted
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Read the raw request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Parse the JSON payload into the GrafanaAlert struct
	var alert GrafanaAlert
	if err := json.Unmarshal(body, &alert); err != nil {
		log.Printf("Failed to unmarshal JSON: %v", err)
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	// Format the alert message using HTML to avoid Telegram Markdown parsing errors (e.g., underscores in rule names)
	message := fmt.Sprintf("🚨 <b>Agritech Alert</b> 🚨\n\n<b>Rule:</b> %s\n<b>State:</b> %s\n<b>Message:</b> %s", 
		alert.RuleName, alert.State, alert.Message)

	// Dispatch the message to the Telegram API
	sendTelegramMessage(botToken, chatID, message)

	// Acknowledge receipt back to Grafana
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Alert received and forwarded successfully."))
}

// sendTelegramMessage pushes the formatted text to the Telegram Bot API
func sendTelegramMessage(botToken string, chatID string, text string) {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken)

	// Prepare the JSON payload for Telegram using HTML parse mode
	payload := map[string]string{
		"chat_id":    chatID,
		"text":       text,
		"parse_mode": "HTML",
	}

	jsonPayload, _ := json.Marshal(payload)

	// Execute the HTTP POST request
	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		log.Printf("Failed to send Telegram message: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Telegram API responded with status: %d", resp.StatusCode)
	} else {
		log.Println("Successfully pushed alert to Telegram.")
	}
}
