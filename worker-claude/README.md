# Worker Claude

- Lắng nghe queue `label-platform-queue-claude` trên Redis
- Gọi mock model Claude, trả về predicted_labels
- Push kết quả vào queue `label-platform-queue-result`

## Biến môi trường
- REDIS_ADDR
- REDIS_PASSWORD 