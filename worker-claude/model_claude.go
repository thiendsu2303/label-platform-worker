package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
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
	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		return PredictionResult{Elements: []UIElement{}}
	}

	prompt := `You are an expert UI analyzer. Given a UI design image in base64 string, your task is to detect and return all UI elements of the following types only: Button, Input, Radio, Drop (Dropdown). For each element, return its type, position (x, y), width, height, and text (if any). For Input, include the placeholder if available.

Output format (JSON): { "elements": [ ... ] }

Only output valid JSON, no explanation.

Here is the base64 image string: ` + base64Image

	reqBody := openrouterRequest{
		Model:    "qwen/qwen3-4b:free",
		Messages: []openMsg{{Role: "user", Content: prompt}},
	}

	bodyBytes, _ := json.Marshal(reqBody)
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", "https://openrouter.ai/api/v1/chat/completions", bytes.NewReader(bodyBytes))
	if err != nil {
		return PredictionResult{Elements: []UIElement{}}
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")
	// req.Header.Set("HTTP-Referer", "<YOUR_SITE_URL>")
	// req.Header.Set("X-Title", "<YOUR_SITE_NAME>")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return PredictionResult{Elements: []UIElement{}}
	}
	defer resp.Body.Close()
	respBytes, _ := io.ReadAll(resp.Body)

	var orResp openrouterResponse
	if err := json.Unmarshal(respBytes, &orResp); err != nil || len(orResp.Choices) == 0 {
		return PredictionResult{Elements: []UIElement{}}
	}

	content := orResp.Choices[0].Message.Content
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
