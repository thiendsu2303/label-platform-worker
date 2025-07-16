package main

import (
	"math/rand"
)

// CallGemini mô phỏng xử lý model Gemini
func CallGemini(base64Image string) PredictionResult {
	return PredictionResult{
		Elements: []UIElement{
			{
				Type:     "dropdown",
				Text:     "Dropdown",
				Width:    150 + rand.Intn(30),
				Height:   35 + rand.Intn(10),
				Position: Position{X: rand.Float64() * 900, Y: rand.Float64() * 450},
			},
			{
				Type:     "slider",
				Width:    200 + rand.Intn(30),
				Height:   20 + rand.Intn(5),
				Position: Position{X: rand.Float64() * 900, Y: rand.Float64() * 450},
			},
		},
	}
}
