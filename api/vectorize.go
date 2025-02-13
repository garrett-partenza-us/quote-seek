package main

import (
	"log"
	"strings"
	"github.com/fogfish/word2vec"
)

type Vectorizer struct {
	Model word2vec.Model
}

func LoadModel(path string) word2vec.Model {
	w2v, err := word2vec.Load("../model/word2vec_model.bin", 265)
	if err != nil {
		log.Fatal(err)
	}
	return w2v
}

func NewVectorizer(path string) Vectorizer {
	model := LoadModel(path)
	return Vectorizer {
		Model: model,
	}
}

func (v *Vectorizer) Vectorize(tokens []string) []float32 {
	vec := make([]float32, 265)
	doc := strings.Join(tokens, " ")
	err := v.Model.Embedding(doc, vec)
	if err != nil {
		log.Fatal(err)
	}
	return vec
}
