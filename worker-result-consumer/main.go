package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

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
		sql := "UPDATE image_samples SET predicted_labels = jsonb_set(coalesce(predicted_labels, '{}'), '{" + msg.Model + "}', ?) WHERE id = ?;"
		if err := db.Exec(sql, string(labelsJson), msg.ImageID).Error; err != nil {
			log.Printf("DB update error: %v", err)
		} else {
			log.Printf("Updated image %s for model %s", msg.ImageID, msg.Model)
		}
	}
}
