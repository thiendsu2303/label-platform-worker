package main

import (
	"math/rand"
)

// CallClaude mô phỏng xử lý model Claude
func CallClaude(base64Image string) PredictionResult {
	return PredictionResult{
		Elements: []UIElement{
			{
				Type:     "checkbox",
				Text:     "Checkbox",
				Width:    60 + rand.Intn(10),
				Height:   30 + rand.Intn(5),
				Position: Position{X: rand.Float64() * 800, Y: rand.Float64() * 400},
			},
			{
				Type:     "label",
				Text:     "Label",
				Width:    120 + rand.Intn(20),
				Height:   30 + rand.Intn(10),
				Position: Position{X: rand.Float64() * 800, Y: rand.Float64() * 400},
			},
		},
	}
}
