package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sort"
)

type Entry struct {
	Text   string
	Vector []float32
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

		text := record[1]

		// Split the vector string (remove the brackets and split by spaces)
		vectorStr := strings.Trim(record[2], "[]")
		vectorParts := strings.Fields(vectorStr)

		// Convert vector parts to float32
		var vector []float32
		for _, part := range vectorParts {
			val, err := strconv.ParseFloat(part, 32)
			if err != nil {
				return nil, fmt.Errorf("error parsing vector value: %v", err)
			}
			vector = append(vector, float32(val))
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

func (m *Meditations) Search (query []float32) []string {

	// Compute distances
	var distances []float32
	for i := range m.Entries {
		distance := ComputeSSD(query, m.Entries[i].Vector)
		distances = append(distances, distance)
	}
	// Retrieve indices of shortest distances
	indices := make([]int, len(distances))
	for i := range distances {
		indices[i] = i
	}
	sort.Slice(indices, func(i, j int) bool {
		return distances[indices[i]] < distances[indices[j]]
	})

	// Retrieve texts from results
	var results []string
	for i := 0; i < 3; i++ {
		results = append(results, m.Entries[indices[i]].Text)
	}

	return results

}

func ComputeSSD(query []float32, vector []float32) float32 {
	var sum float32
	for i := range query{
		diff := query[i] - vector[i]
		sum += diff * diff
	}
	return sum
}
