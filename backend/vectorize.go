package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"log"
	"os"
	"io/ioutil"
	"net/http"
	"fmt"
)

type Vectorizer struct {
	Vocab   map[string]int
	Vectors [][]float64
	Ngrams  [][]float64
	Bucket int
}

func LoadVocab(path string) map[string]int {
	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	var vocab map[string]int

	err = json.Unmarshal(file, &vocab)
	if err != nil {
		log.Fatal(err)
	}

	return vocab
}

func ReadBinaryArrayFile2D(path string, rows int, cols int) [][]float64 {

	// Read the binary array file
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Get the file size in bytes
	fileInfo, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}
	fileSize := fileInfo.Size()

	// Ensure file bytes is multiple of float64 byte size
	if fileSize%8 != 0 {
		log.Fatal("The file size is not a multiple of 4 bytes (size of float64)")
	}

	// Read the binary data into a byte slice
	data := make([]byte, fileSize)
	_, err = file.Read(data)
	if err != nil {
		log.Fatal(err)
	}

	// Convert the byte slice into a slice of float64
	numFloats := fileSize / 8
	floats := make([]float64, numFloats)
	err = binary.Read(bytes.NewReader(data), binary.LittleEndian, &floats)
	if err != nil {
		log.Fatal(err)
	}

	// Check if the number of elements matches the expected 2D array size
	if rows*cols != len(floats) {
		log.Fatal("The dimensions of the array do not match the expected size.")
	}

	// Create a 2D slice and fill it with the data from the flat slice
	array2D := make([][]float64, rows)
	for i := range array2D {
		array2D[i] = floats[i*cols : (i+1)*cols]
	}

	return array2D
}

func NewVectorizer(vocabPath string, vectorsPath string, ngramsPath string, size int, bucket int) Vectorizer {

	vocab := LoadVocab(vocabPath)
	vectors := ReadBinaryArrayFile2D(vectorsPath, len(vocab), size)
	ngrams := ReadBinaryArrayFile2D(ngramsPath, bucket, size)

	return Vectorizer{
		Vocab:   vocab,
		Vectors: vectors,
		Ngrams:  ngrams,
		Bucket: bucket,
	}
}

func customFTHashBytes(bytez []byte) uint32 {
	h := uint32(2166136261) // Initial FNV offset basis
	for _, b := range bytez {
		h ^= uint32(b) // XOR the current byte value
		h *= 16777619  // Multiply by the FNV prime
	}
	return h
}

func generateCharNGrams(text string, n int) []string {
	if len(text) < n || n <= 0 {
		return []string{}
	}
	var ngrams []string
	for i := 0; i <= len(text)-n; i++ {
		ngrams = append(ngrams, text[i:i+n])
	}
	return ngrams
}

func generateNGrams(word string, minN int, maxN int) []string {
	var ngrams []string
	for n := minN; n <= maxN; n++ {
		ngrams = append(ngrams, generateCharNGrams(word, n)...)
	}
	return ngrams
}

func addSlices(a []float64, b []float64) []float64 {
	if len(a) != len(b) {
		log.Fatal("Cannot add slices element wise of differing lengths.")
	}

	result := make([]float64, len(a))
	for i := 0; i < len(a); i++ {
		result[i] = a[i] + b[i]
	}

	return result
}

func averageSlice(arr []float64, d int) []float64 {
	result := make([]float64, len(arr))

	for i := 0; i < len(arr); i++ {
		result[i] = arr[i] / float64(d)
	}

	return result
}

func (v *Vectorizer) Vectorize(tokens []string) []float64 {

	vectors := make([][]float64, len(tokens))

	for idx, token := range tokens {
		if _, exists := v.Vocab[token]; exists {
			vectors[idx] = v.Vectors[v.Vocab[token]]
		} else {
			ngrams := generateNGrams("<"+token+">", 3, 6)
			vector := make([]float64, 256)
			for _, ngram := range ngrams {
				utf8Bytes := []byte(ngram)
				boundedHash := customFTHashBytes(utf8Bytes) % uint32(v.Bucket)
				ngramVector := v.Ngrams[boundedHash]
				vector = addSlices(vector, ngramVector)
			}
			vectors[idx] = averageSlice(vector, len(ngrams))
		}
	}

	result := make([]float64, 256)
	for _, vector := range vectors {
		result = addSlices(result, vector)
	}
	return averageSlice(result, len(vectors))
}

type EmbeddingRequest struct {
	Input string `json:"input"`
	Model string `json:"model"`
}

type EmbeddingResponse struct {
	Object string `json:"object"`
	Data   []struct {
		Object   string    `json:"object"`
		Index    int       `json:"index"`
		Embedding []float64 `json:"embedding"`
	} `json:"data"`
	Model string `json:"model"`
	Usage struct {
		PromptTokens  int `json:"prompt_tokens"`
		TotalTokens   int `json:"total_tokens"`
	} `json:"usage"`
}

func getEmbedding(text string) ([]float64, error) {

	apiKey := os.Getenv("openai_key")
	if apiKey == "" {
		return nil, fmt.Errorf("API key is missing")
	}

	url := "https://api.openai.com/v1/embeddings"

	embeddingRequest := EmbeddingRequest{
		Input: text,
		Model: "text-embedding-3-small",
	}
	requestBody, err := json.Marshal(embeddingRequest)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making API request: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status code %d: %s", resp.StatusCode, string(respBody))
	}

	var embeddingResponse EmbeddingResponse
	if err := json.Unmarshal(respBody, &embeddingResponse); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	if len(embeddingResponse.Data) > 0 {
		return embeddingResponse.Data[0].Embedding, nil
	}
	return nil, fmt.Errorf("no embedding data returned")
}
