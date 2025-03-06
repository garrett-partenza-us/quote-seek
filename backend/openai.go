package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type ChatGPT struct {
	OpenAI_Environment  string
	OpenAI_Model        string
	OpenAI_Key          string
	OpenAI_Endpoint     string
	OpenAI_SystemPrompt string
	OpenAI_MaxTokens    int
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ResponseFormat struct {
	Type       string     `json:"type"`
	JsonSchema JsonSchema `json:"json_schema"`
}

type Type struct {
	Type string `json:"type"`
}

type Properties struct {
	Quote          Type `json:"quote"`
	Interpretation Type `json:"interpretation"`
	Advice         Type `json:"advice"`
}

type Schema struct {
	Type                 string     `json:"type"`
	Properties           Properties `json:"properties"`
	Required             []string   `json:"required"`
	AdditionalProperties bool       `json:"additionalProperties"`
}

type JsonSchema struct {
	Name   string `json:"name"`
	Schema Schema `json:"schema"`
	Strict bool   `json:"strict"`
}

type RequestData struct {
	Model               string         `json:"model"`
	Messages            []Message      `json:"messages"`
	ResponseFormat      ResponseFormat `json:"response_format"`
	Store               bool           `json:"store"`
	MaxCompletionTokens int            `json:"max_completion_tokens"`
	N                   int            `json:"n"`
	User                string         `json:"user"`
}

type Choice struct {
	Index   int     `json:"index"`
	Message Message `json:"message"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type APIResponse struct {
	ID          string   `json:"id"`
	Object      string   `json:"object"`
	Created     int64    `json:"created"`
	Model       string   `json:"model"`
	Choices     []Choice `json:"choices"`
	Usage       Usage    `json:"usage"`
	ServiceTier string   `json:"service_tier"`
}
// Invalid schema for response_format 'model_output': None is not of type 'array'.

func (c *ChatGPT) Query(prompt string) (string, string, string) {

	data := RequestData{
		Model: c.OpenAI_Model,
		Messages: []Message{
			{
				Role:    "system",
				Content: c.OpenAI_SystemPrompt,
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
		ResponseFormat: ResponseFormat{
			Type: "json_schema",
			JsonSchema: JsonSchema{
				Name: "model_output",
				Schema: Schema{
					Type: "object",
					Properties: Properties{
						Quote: Type{
							Type: "string",
						},
						Interpretation: Type{
							Type: "string",
						},
						Advice: Type{
							Type: "string",
						},
					},
					Required: []string{"quote", "interpretation", "advice"},
					AdditionalProperties: false,
				},
				Strict: true,
			},
		},
		Store:               false,
		MaxCompletionTokens: c.OpenAI_MaxTokens,
		N:                   1,
		User:                c.OpenAI_Environment,
	}

	payload, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("Error marshalling data: %v", err)
	}

	client := &http.Client{}

	req, err := http.NewRequest("POST", c.OpenAI_Endpoint, bytes.NewBuffer(payload))
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+c.OpenAI_Key)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response: %v", err)
	}

	fmt.Println(string(body))

	var response APIResponse
	err = json.Unmarshal([]byte(body), &response)
	if err != nil {
		log.Fatalf("Error unmarshalling response: %v", err)
	}

	var content map[string]string
	err = json.Unmarshal([]byte(response.Choices[0].Message.Content), &content)
	if err != nil {
		log.Fatalf("Error unmarshalling content: %v", err)
	}

	fmt.Println("Quote:", content["quote"])
	fmt.Println("Interpretation:", content["interpretation"])
	fmt.Println("Advice:", content["advice"])

	log.Printf("Prompt Tokens: %d\n", response.Usage.PromptTokens)
	log.Printf("Completion Tokens: %d\n", response.Usage.CompletionTokens)
	log.Printf("Total Tokens: %d\n", response.Usage.TotalTokens)

	return content["quote"], content["interpretation"], content["advice"]

}
