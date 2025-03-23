package main

import (
	"encoding/csv"
	"fmt"
	"math"
	"io"
	"os"
	"strconv"
	"strings"
	"sort"
)

type Entry struct {
	Text   string
	Vector []float64
}

type Meditations struct {
	NumEntries int
	Entries    []Entry
}

func NewMeditations(path string) (*Meditations, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Read header, skip it
	_, err = reader.Read()
	if err != nil {
		return nil, err
	}

	var entries []Entry

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		// record[1] is the chunk (text)
		// record[2] is the vector (embeddings)

		text := record[0]

		// Split the vector string (remove the brackets and split by spaces)
		vectorStr := strings.Trim(record[1], "[]")
		vectorStr = strings.ReplaceAll(vectorStr, ",", "") // Remove commas
		vectorParts := strings.Fields(vectorStr)

		// Convert vector parts to float64
		var vector []float64
		for _, part := range vectorParts {
			val, err := strconv.ParseFloat(part, 64)
			if err != nil {
				return nil, fmt.Errorf("error parsing vector value: %v", err)
			}
			vector = append(vector, float64(val))
		}

		entries = append(entries, Entry{
			Text:   text,
			Vector: vector,
		})
	}

	return &Meditations{
		NumEntries: len(entries),
		Entries:    entries,
	}, nil
}

func (m *Meditations) Search (query []float64) []string {

	// Compute distances
	var distances []float64
	for i := range m.Entries {
		distance := cosineSimilarity(query, m.Entries[i].Vector)
		distances = append(distances, distance)
	}
	// Retrieve indices of shortest distances
	indices := make([]int, len(distances))
	for i := range distances {
		indices[i] = i
	}
	sort.Slice(indices, func(i, j int) bool {
		return distances[indices[i]] > distances[indices[j]]
	})

	// Retrieve texts from results
	var results []string
	for i := 0; i < 5; i++ {
		results = append(results, m.Entries[indices[i]].Text)
	}

	return results

}

// Function to compute the dot product of two vectors
func dotProduct(a, b []float64) float64 {
	var result float64
	for i := 0; i < len(a); i++ {
		result += a[i] * b[i]
	}
	return result
}

// Function to compute the magnitude (Euclidean norm) of a vector
func magnitude(a []float64) float64 {
	var sum float64
	for _, value := range a {
		sum += value * value
	}
	return float64(math.Sqrt(float64(sum)))
}

// Function to compute the cosine similarity
func cosineSimilarity(a, b []float64) float64 {
	// Calculate the dot product of vectors a and b
	dot := dotProduct(a, b)
	// Calculate the magnitudes of vectors a and b
	magnitudeA := magnitude(a)
	magnitudeB := magnitude(b)

	// Return the cosine similarity
	return dot / (magnitudeA * magnitudeB)
}
