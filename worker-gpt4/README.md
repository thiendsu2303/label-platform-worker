# Worker GPT-4

- Lắng nghe queue `label-platform-queue-gpt` trên Redis
- Gọi mock model GPT-4, trả về predicted_labels
- Push kết quả vào queue `label-platform-queue-result`

## Biến môi trường
- REDIS_ADDR
- REDIS_PASSWORD 