package main

import (
	"fmt"
	"encoding/json"
	"log"
	"io/ioutil"
	"net/http"
)

type Request struct {
	Query string
}

func Search(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST"{
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

	fmt.Println(data.Query)

	// Clean text
	// Vectorize
	// Search
	// Chat

}

func main() {
	config := GetConfig()
	normalizer := NewNormalizer()
	tokens := normalizer.Normalize("universe plants stars sky")
	vectorizer := NewVectorizer("../model/word2vec_model")
	vector := vectorizer.Vectorize(tokens)
	fmt.Println(vector)
	meditations, err := NewMeditations(config.MeditationsCSVPath)
	result := meditations.Search(vector)
	for _, res := range result {
		fmt.Println(res)
	}
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/search", Search)

	fmt.Println("Server is listening on port 8080...")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
