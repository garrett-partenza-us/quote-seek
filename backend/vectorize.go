package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"log"
	"os"
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
