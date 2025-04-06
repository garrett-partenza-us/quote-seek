package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

// Define a struct for feedback data
type FeedbackData struct {
	Name       string `json:"name"`
	Consistent string `json:"consistent"`
	Helpful    string `json:"helpful"`
}

type SearchHandler struct {
	Meditations *Meditations
	ChatGPT     ChatGPT
}

type Request struct {
	Query string
}

// FeedbackHandler handles feedback submissions
type FeedbackHandler struct {
	DB *sql.DB
}

// Enable CORS for all origins and required methods
func enableCors(w *http.ResponseWriter, r *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == "OPTIONS" {
		(*w).WriteHeader(http.StatusOK)
		return
	}
}

func (h SearchHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	enableCors(&w, r)
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusInternalServerError)
		return
	}
	var data Request
	err = json.Unmarshal(body, &data)
	if err != nil {
		http.Error(w, "Failed to unmarshal data", http.StatusInternalServerError)
		return
	}
	//cleanedText := h.Normalizer.Normalize(data.Query)
	//queryVector := h.Vectorizer.Vectorize(cleanedText)
	//scaledVector := h.StandardScaler.ScaleVector(queryVector)
	scaledVector, err := getEmbedding(data.Query)
	if err != nil {
		log.Fatal(err)
	}
	searchResults := h.Meditations.Search(scaledVector)
	userPrompt := GeneratePrompt(data.Query, searchResults)
	quote, interpretation, advice := h.ChatGPT.Query(userPrompt)
	response := map[string]string{
		"quote":          quote,
		"interpretation": interpretation,
		"advice":         advice,
	}
	w.Header().Set("Content-Type", "application/json")
	respJSON, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(respJSON)
}

// ServeHTTP handles HTTP requests for feedback submissions
func (h FeedbackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	enableCors(&w, r)

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusInternalServerError)
		return
	}

	// Parse JSON data
	var feedbackData FeedbackData
	err = json.Unmarshal(body, &feedbackData)
	if err != nil {
		http.Error(w, "Failed to unmarshal data", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if feedbackData.Name == "" || feedbackData.Consistent == "" || feedbackData.Helpful == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// Insert data into database
	_, err = h.DB.Exec(
		"INSERT INTO feedback (name, consistent, helpful) VALUES (?, ?, ?)",
		feedbackData.Name,
		feedbackData.Consistent,
		feedbackData.Helpful,
	)

	if err != nil {
		log.Printf("Database error: %v", err)
		http.Error(w, "Failed to store feedback", http.StatusInternalServerError)
		return
	}

	// Return success response
	response := map[string]string{
		"status":  "success",
		"message": "Feedback submitted successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	respJSON, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(respJSON)
}

// initDB initializes the database connection and creates the table if it doesn't exist
func initDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Create the feedback table if it doesn't exist
	createTableSQL := `
    CREATE TABLE IF NOT EXISTS feedback (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        consistent TEXT NOT NULL,
        helpful TEXT NOT NULL,
        submitted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );
    
    CREATE INDEX IF NOT EXISTS idx_feedback_date ON feedback(submitted_at);
    `

	_, err = db.Exec(createTableSQL)
	if err != nil {
		return nil, fmt.Errorf("failed to create table: %v", err)
	}

	return db, nil
}

func main() {
	// Initialize database connection
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		// Use default path if not specified
		dbPath = "./db/feedback.db"
	}

	db, err := initDB(dbPath)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	meditations, err := NewMeditations(os.Getenv("meditations_csv_path"))
	if err != nil {
		log.Fatal("Failed to load meditations: ", err)
	}
	maxTokens, err := strconv.Atoi(os.Getenv("openai_maxtokens")) // TODO: Refactor this later
	if err != nil {
		log.Fatal(err)
	}
	chatgpt := ChatGPT{
		OpenAI_Environment:  os.Getenv("openai_environment"),
		OpenAI_Model:        os.Getenv("openai_model"),
		OpenAI_Key:          os.Getenv("openai_key"),
		OpenAI_Endpoint:     os.Getenv("openai_endpoint"),
		OpenAI_SystemPrompt: os.Getenv("openai_systemprompt"),
		OpenAI_MaxTokens:    maxTokens,
	}

	// Create handlers
	searchHandler := SearchHandler{
		Meditations: meditations,
		ChatGPT:     chatgpt,
	}

	feedbackHandler := FeedbackHandler{
		DB: db,
	}

	// Set up routes
	http.Handle("/search", searchHandler)
	http.Handle("/feedback", feedbackHandler)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server is listening on port %s...\n", port)
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
