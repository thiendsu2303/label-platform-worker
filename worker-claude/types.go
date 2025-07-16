package main

type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type UIElement struct {
	Type        string   `json:"type"`
	Text        string   `json:"text,omitempty"`
	Width       int      `json:"width"`
	Height      int      `json:"height"`
	Position    Position `json:"position"`
	Placeholder string   `json:"placeholder,omitempty"`
}

type JobPayload struct {
	ImageID     string `json:"image_id"`
	Base64Image string `json:"base64_image"`
}

type PredictionResult struct {
	Elements []UIElement `json:"elements"`
}

type PredictionResultMessage struct {
	ImageID         string           `json:"image_id"`
	Model           string           `json:"model"`
	Base64Image     string           `json:"base64_image"`
	PredictedLabels PredictionResult `json:"predicted_labels"`
}
