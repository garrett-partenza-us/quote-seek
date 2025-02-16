package main

import (
	"strings"
	"fmt"
)

func GeneratePrompt(query string, searchResults []string, preface string) string {

	// Construct query XML
	queryXML := "<rag_query>" + query + "</rag_query>"

	// Construct search results XML
	var sb strings.Builder
	sb.WriteString("<search_results>")
	for _, result := range searchResults {
		sb.WriteString(fmt.Sprintf("<search_result>%s</search_result>", result))
	}
	sb.WriteString("</search_results>")
	searchResultsXML := sb.String()

	// Create user prompt
	sb.Reset()
	sb.WriteString("<user_prompt>")
	sb.WriteString(preface)
	sb.WriteString(queryXML)
	sb.WriteString(searchResultsXML)
	sb.WriteString("</user_prompt>")
	userPrompt := sb.String()

	return userPrompt

}
