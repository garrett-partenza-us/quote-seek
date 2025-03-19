package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

type SearchHandler struct {
	Meditations					*Meditations
	Normalizer					Normalizer
	Vectorizer					Vectorizer
	StandardScaler			StandardScaler
	ChatGPT							ChatGPT
}

type Request struct {
	Query string
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

	cleanedText := h.Normalizer.Normalize(data.Query)
	queryVector := h.Vectorizer.Vectorize(cleanedText)
	scaledVector := h.StandardScaler.ScaleVector(queryVector)
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

func main() {
	meditations, err := NewMeditations(os.Getenv("meditations_csv_path"))
	if err != nil {
		log.Fatal("Failed to load meditations: ", err)
	}
	normalizer := NewNormalizer(os.Getenv("stopwords_path"))
	vectorizer := NewVectorizer(
		os.Getenv("vocab_path"),
		os.Getenv("vectors_path"),
		os.Getenv("ngrams_path"),
		256,
		100000,
	)

	scaler := NewStandardScaler(
		256,
		os.Getenv("mean_path"),
		os.Getenv("scale_path"),
	)

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
	handler := SearchHandler{
		Meditations:			meditations,
		Normalizer:				normalizer,
		Vectorizer:				vectorizer,
		ChatGPT:					chatgpt,
		StandardScaler:		scaler,
	}
	if err != nil {
		log.Fatal(err)
	}
	http.Handle("/search", handler)
	fmt.Println("Server is listening on port 8080...")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
