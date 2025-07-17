package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type openrouterRequest struct {
	Model    string    `json:"model"`
	Messages []openMsg `json:"messages"`
}

type openMsg struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openrouterResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func CallClaude(base64Image string) PredictionResult {
	apiKey := os.Getenv("CLAUDE_API_KEY")
	if apiKey == "" {
		log.Printf("CallClaude: missing CLAUDE_API_KEY")
		return PredictionResult{Elements: []UIElement{}}
	}

	prompt := `You are an expert UI analyzer. Given a UI design image in base64 string, your task is to detect and return all UI elements of the following types only: Button, Input, Radio, Drop (Dropdown). For each element, return its type, position (x, y), width, height, and text (if any). For Input, include the placeholder if available.\n\nOutput format (JSON): { \"elements\": [ ... ] }\n\nOnly output valid JSON, no explanation.\n\nHere is the base64 image string: ` + base64Image

	reqBody := openrouterRequest{
		Model:    "qwen/qwen3-4b:free",
		Messages: []openMsg{{Role: "user", Content: prompt}},
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		log.Printf("CallClaude: json.Marshal error: %v", err)
		return PredictionResult{Elements: []UIElement{}}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", "https://openrouter.ai/api/v1/chat/completions", bytes.NewReader(bodyBytes))
	if err != nil {
		log.Printf("CallClaude: NewRequestWithContext error: %v", err)
		return PredictionResult{Elements: []UIElement{}}
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")
	// req.Header.Set("HTTP-Referer", "<YOUR_SITE_URL>")
	// req.Header.Set("X-Title", "<YOUR_SITE_NAME>")

	log.Printf("Call OpenRouter Claude API with prompt length: %d", len(prompt))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("OpenRouter API error: %v", err)
		return PredictionResult{Elements: []UIElement{}}
	}
	defer resp.Body.Close()
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("CallClaude: io.ReadAll error: %v", err)
		return PredictionResult{Elements: []UIElement{}}
	}

	var orResp openrouterResponse
	if err := json.Unmarshal(respBytes, &orResp); err != nil || len(orResp.Choices) == 0 {
		log.Printf("OpenRouter API unmarshal error: %v, body: %s", err, string(respBytes))
		return PredictionResult{Elements: []UIElement{}}
	}

	content := orResp.Choices[0].Message.Content
	log.Printf("OpenRouter API response: %s", content)
	start := strings.Index(content, "{")
	end := strings.LastIndex(content, "}")
	if start == -1 || end == -1 || end <= start {
		log.Printf("CallClaude: cannot find JSON in response: %s", content)
		return PredictionResult{Elements: []UIElement{}}
	}
	jsonStr := content[start : end+1]
	var result PredictionResult
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		log.Printf("CallClaude: json.Unmarshal error: %v, json: %s", err, jsonStr)
		return PredictionResult{Elements: []UIElement{}}
	}
	return result
}
