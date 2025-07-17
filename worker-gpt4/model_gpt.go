package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

// CallGPT4 gọi OpenAI GPT-4 thường với prompt chi tiết, trả về PredictionResult đúng format ground_truth
func CallGPT4(base64Image string) PredictionResult {
	apiKey := os.Getenv("OPENAI_API_KEY")

	client := openai.NewClient(apiKey)
	prompt := `You are an expert UI analyzer. Given a UI design image in base64 string, your task is to detect and return all UI elements of the following types only: Button, Input, Radio, Drop (Dropdown). For each element, return its type, position (x, y), width, height, and text (if any). For Input, include the placeholder if available.

Output format (JSON): { "elements": [ ... ] }

Only output valid JSON, no explanation.

Here is the base64 image string: ` + base64Image

	log.Printf("Call OpenAI GPT-4 API with prompt length: %d", len(prompt))
	resp, err := client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
		Model: "gpt-3.5-turbo",
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		MaxTokens:   2048,
		Temperature: 0,
	})
	if err != nil || len(resp.Choices) == 0 {
		log.Printf("OpenAI API error: %v", err)
		return PredictionResult{Elements: []UIElement{}}
	}
	content := resp.Choices[0].Message.Content
	log.Printf("OpenAI API response: %s", content)
	// Tìm đoạn JSON trong content (nếu có text thừa)
	start := strings.Index(content, "{")
	end := strings.LastIndex(content, "}")
	if start == -1 || end == -1 || end <= start {
		return PredictionResult{Elements: []UIElement{}}
	}
	jsonStr := content[start : end+1]
	var result PredictionResult
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		return PredictionResult{Elements: []UIElement{}}
	}
	return result
}
