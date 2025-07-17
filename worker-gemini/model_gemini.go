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

type geminiRequest struct {
	Contents []struct {
		Parts []struct {
			Text string `json:"text"`
		} `json:"parts"`
	} `json:"contents"`
}

type geminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

func CallGemini(base64Image string) PredictionResult {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Printf("CallGemini: missing GEMINI_API_KEY")
		return PredictionResult{Elements: []UIElement{}}
	}

	// Prompt truyền base64 image như cũ
	prompt := `You are an expert UI analyzer. Given a UI design image in base64 string, your task is to detect and return all UI elements of the following types only: Button, Input, Radio, Drop (Dropdown). For each element, return its type, position (x, y), width, height, and text (if any). For Input, include the placeholder if available.

Output format (JSON): { "elements": [ ... ] }

Only output valid JSON, no explanation.

Here is the base64 image string: ` + base64Image

	reqBody := geminiRequest{
		Contents: []struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		}{
			{
				Parts: []struct {
					Text string `json:"text"`
				}{
					{Text: prompt},
				},
			},
		},
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		log.Printf("CallGemini: json.Marshal error: %v", err)
		return PredictionResult{Elements: []UIElement{}}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	url := "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.5-flash:generateContent"
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(bodyBytes))
	if err != nil {
		log.Printf("CallGemini: NewRequestWithContext error: %v", err)
		return PredictionResult{Elements: []UIElement{}}
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-goog-api-key", apiKey)

	log.Printf("Call Gemini API (REST) with prompt length: %d", len(prompt))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Gemini API error: %v", err)
		return PredictionResult{Elements: []UIElement{}}
	}
	defer resp.Body.Close()
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("CallGemini: io.ReadAll error: %v", err)
		return PredictionResult{Elements: []UIElement{}}
	}

	var gResp geminiResponse
	if err := json.Unmarshal(respBytes, &gResp); err != nil || len(gResp.Candidates) == 0 {
		log.Printf("Gemini API unmarshal error: %v, body: %s", err, string(respBytes))
		return PredictionResult{Elements: []UIElement{}}
	}

	var content string
	for _, part := range gResp.Candidates[0].Content.Parts {
		content += part.Text
	}
	log.Printf("Gemini API response: %s", content)
	start := strings.Index(content, "{")
	end := strings.LastIndex(content, "}")
	if start == -1 || end == -1 || end <= start {
		log.Printf("CallGemini: cannot find JSON in response: %s", content)
		return PredictionResult{Elements: []UIElement{}}
	}
	jsonStr := content[start : end+1]
	var result PredictionResult
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		log.Printf("CallGemini: json.Unmarshal error: %v, json: %s", err, jsonStr)
		return PredictionResult{Elements: []UIElement{}}
	}
	return result
}
