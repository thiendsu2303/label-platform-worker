package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

func main() {
	redisAddr := os.Getenv("REDIS_ADDR")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisDB := 0
	queueName := "label-platform-queue-gpt"
	resultQueue := "label-platform-queue-result"
	modelName := "gpt-4"

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
		var job JobPayload
		if err := json.Unmarshal([]byte(result[1]), &job); err != nil {
			log.Printf("Invalid job payload: %v", err)
			continue
		}
		resultObj := CallGPT4(job.Base64Image)
		msg := PredictionResultMessage{
			ImageID:         job.ImageID,
			Model:           modelName,
			Base64Image:     job.Base64Image,
			PredictedLabels: resultObj,
		}
		msgBytes, _ := json.Marshal(msg)
		if err := rdb.RPush(ctx, resultQueue, string(msgBytes)).Err(); err != nil {
			log.Printf("Push result error: %v", err)
		} else {
			log.Printf("Processed image %s", job.ImageID)
		}
	}
}
