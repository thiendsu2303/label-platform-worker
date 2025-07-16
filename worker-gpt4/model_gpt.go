package main

import (
	"math/rand"
)

// CallGPT4 mô phỏng xử lý model GPT-4
func CallGPT4(base64Image string) PredictionResult {
	return PredictionResult{
		Elements: []UIElement{
			{
				Type:     "button",
				Text:     "Button",
				Width:    100 + rand.Intn(20),
				Height:   40 + rand.Intn(10),
				Position: Position{X: rand.Float64() * 1000, Y: rand.Float64() * 500},
			},
			{
				Type:        "input",
				Width:       280 + rand.Intn(20),
				Height:      50 + rand.Intn(10),
				Position:    Position{X: rand.Float64() * 1000, Y: rand.Float64() * 500},
				Placeholder: "Input",
			},
		},
	}
}
