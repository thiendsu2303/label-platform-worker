package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func notifyBackend(imageID, model string, result any) {
	backendURL := os.Getenv("BACKEND_NOTIFY_URL")

	payload := map[string]any{
		"image_id": imageID,
		"model":    model,
		"result":   "success",
	}
	body, _ := json.Marshal(payload)
	log.Printf("Notify backend: POST %s\nPayload: %s", backendURL, string(body))
	req, err := http.NewRequest("POST", backendURL, bytes.NewReader(body))
	if err != nil {
		log.Printf("Notify backend error (create request): %v", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Notify backend error (do request): %v", err)
		return
	}
	defer resp.Body.Close()
	log.Printf("Notify backend response status: %s", resp.Status)
	respBody, _ := io.ReadAll(resp.Body)
	log.Printf("Notify backend response body: %s", string(respBody))
}

func main() {
	redisAddr := os.Getenv("REDIS_ADDR")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisDB := 0
	queueName := "label-platform-queue-result"
	dsn := os.Getenv("POSTGRES_DSN")
	if dsn == "" {
		log.Fatal("POSTGRES_DSN is required")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       redisDB,
	})
	ctx := context.Background()

	for {
		result, err := rdb.BLPop(ctx, 0, queueName).Result()
		if err != nil || len(result) < 2 {
			log.Printf("BLPop error: %v", err)
			time.Sleep(time.Second)
			continue
		}
		var msg PredictionResultMessage
		if err := json.Unmarshal([]byte(result[1]), &msg); err != nil {
			log.Printf("Invalid result message: %v", err)
			continue
		}
		labelsJson, _ := json.Marshal(msg.PredictedLabels)
		sql := "UPDATE images SET predicted_labels = jsonb_set(coalesce(predicted_labels, '{}'), '{" + msg.Model + "}', ?) WHERE id = ?;"
		if err := db.Exec(sql, string(labelsJson), msg.ImageID).Error; err != nil {
			log.Printf("DB update error: %v", err)
		} else {
			log.Printf("Updated image %s for model %s", msg.ImageID, msg.Model)
			notifyBackend(msg.ImageID, msg.Model, msg.PredictedLabels)
		}
	}
}
