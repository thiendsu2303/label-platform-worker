package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"log"
	"os"
	"strings"

	"cloud.google.com/go/vertexai/genai"
	"google.golang.org/api/option"
)

func CallGemini(base64Image string) PredictionResult {
	apiKey := os.Getenv("GEMINI_API_KEY")
	projectID := os.Getenv("GEMINI_PROJECT_ID")
	location := os.Getenv("GEMINI_LOCATION")
	if apiKey == "" || projectID == "" || location == "" {
		return PredictionResult{Elements: []UIElement{}}
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, projectID, location, option.WithAPIKey(apiKey))
	if err != nil {
		return PredictionResult{Elements: []UIElement{}}
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-pro-vision")

	// Decode base64 image to []byte
	imgBytes, err := base64.StdEncoding.DecodeString(base64Image)
	if err != nil {
		return PredictionResult{Elements: []UIElement{}}
	}

	prompt := `You are an expert UI analyzer. Given a UI design image, your task is to detect and return all UI elements of the following types only: Button, Input, Radio, Drop (Dropdown). For each element, return its type, position (x, y), width, height, and text (if any). For Input, include the placeholder if available.

Output format (JSON): { "elements": [ ... ] }

Only output valid JSON, no explanation.`

	// Gửi prompt và ảnh lên Gemini
	resp, err := model.GenerateContent(ctx, genai.Text(prompt), genai.ImageData("image/png", imgBytes))
	if err != nil || len(resp.Candidates) == 0 {
		return PredictionResult{Elements: []UIElement{}}
	}

	// Lấy text từ các part
	var content string
	for _, part := range resp.Candidates[0].Content.Parts {
		if text, ok := part.(*genai.Text); ok {
			content += string(*text)
		}
	}
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
	log.Printf("Model raw response: %s", content)
	return result
}
