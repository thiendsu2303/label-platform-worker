# Worker Gemini

- Lắng nghe queue `label-platform-queue-gemini` trên Redis
- Gọi mock model Gemini, trả về predicted_labels
- Push kết quả vào queue `label-platform-queue-result`

## Biến môi trường
- REDIS_ADDR
- REDIS_PASSWORD 