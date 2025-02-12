package main

import (
	"bufio"
	"log"
	"os"
	"slices"
	"strings"
	"unicode"
)

type Normalizer struct {
	Stopwords []string
}

func NewNormalizer() Normalizer {
	stopwords := LoadStopwords("./stopwords-en.txt")
	return Normalizer{
		Stopwords: stopwords,
	}
}

func (n *Normalizer) Normalize(text string) []string {
	text = n.Lower(text)
	text = n.RemovePunctuation(text)
	tokens := n.Tokenize(text)
	tokens = n.RemoveStopwords(tokens)
	return tokens
}

func LoadStopwords(path string) []string {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var stopwords []string

	for scanner.Scan() {
		stopwords = append(stopwords, strings.TrimSpace(scanner.Text()))
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return stopwords
}

func (n *Normalizer) Lower(text string) string {
	return strings.ToLower(text)
}

func (n *Normalizer) RemovePunctuation(text string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsPunct(r) {
			return -1
		}
		return r
	}, text)
}

func (n *Normalizer) Tokenize(text string) []string {
	return strings.Fields(text)
}

func (n *Normalizer) RemoveStopwords(tokens []string) []string {
	var result []string
	for _, token := range tokens {
		if !(slices.Contains(n.Stopwords, token)) {
			result = append(result, token)
		}
	}
	return result
}
