package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type SearchHandler struct {
	Config      Config
	Meditations *Meditations
	Normalizer  Normalizer
	Vectorizer  Vectorizer
	ChatGPT     ChatGPT
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
	searchResults := h.Meditations.Search(queryVector)

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
	config := GetConfig()
	meditations, err := NewMeditations(config.MeditationsCSVPath)
	normalizer := NewNormalizer(config.StopwordsPath)
	vectorizer := NewVectorizer(
		config.VocabPath,
		config.VectorsPath,
		config.NgramsPath,
		256,
		100000,
	)
	chatgpt := ChatGPT{
		OpenAI_Environment:  config.OpenAI_Environment,
		OpenAI_Model:        config.OpenAI_Model,
		OpenAI_Key:          config.OpenAI_Key,
		OpenAI_Endpoint:     config.OpenAI_Endpoint,
		OpenAI_SystemPrompt: config.OpenAI_SystemPrompt,
		OpenAI_MaxTokens:    config.OpenAI_MaxTokens,
	}
	handler := SearchHandler{
		Config:      config,
		Meditations: meditations,
		Normalizer:  normalizer,
		Vectorizer:  vectorizer,
		ChatGPT:     chatgpt,
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
