package main

import (
	"bytes"
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
}

type Request struct {
	Query string
}

func (h SearchHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

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

	userPrompt := GeneratePrompt(data.Query, searchResults, h.Config.Preface)

	// Prepare the JSON payload to send to the Ollama API
	payload := map[string]interface{}{
		"model": "llama3:instruct",
		"messages": []map[string]interface{}{
			{
				"role":    "system",
				"content": h.Config.SystemPrompt,
			},
			{
				"role":    "user",
				"content": userPrompt,
			},
		},
		"stream": false,
	}

	ollamaURL := "http://localhost:11434/api/chat"
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("POST", ollamaURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to call Ollama API", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Read response from Ollama
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response from Ollama", http.StatusInternalServerError)
		return
	}

	var ollamaResponse map[string]interface{}
	err = json.Unmarshal(respBody, &ollamaResponse)
	if err != nil {
		http.Error(w, "Failed to parse Ollama response", http.StatusInternalServerError)
		return
	}

	// Ensure "message" and "content" keys exist before accessing them
	message, ok := ollamaResponse["message"].(map[string]interface{})
	if !ok {
		http.Error(w, "Invalid response format: missing 'message' field", http.StatusInternalServerError)
		return
	}

	content, ok := message["content"].(string)
	if !ok {
		http.Error(w, "Invalid response format: missing 'content' field", http.StatusInternalServerError)
		return
	}

	// Send the generated text back to the original caller
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	encoder := json.NewEncoder(w)
	response := map[string]string{"response": content}
	encoder.Encode(response)
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
	handler := SearchHandler{
		Config:      config,
		Meditations: meditations,
		Normalizer:  normalizer,
		Vectorizer:  vectorizer,
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
