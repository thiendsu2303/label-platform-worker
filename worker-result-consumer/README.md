# Worker Result Consumer

- Lắng nghe queue `label-platform-queue-result` trên Redis
- Cập nhật trường predicted_labels.<model> vào bảng image_samples trên PostgreSQL

## Biến môi trường
- REDIS_ADDR
- REDIS_PASSWORD
- POSTGRES_DSN 