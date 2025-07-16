# Label Platform Worker

Hệ thống gồm 4 service Golang chạy độc lập (Docker container), xử lý gán nhãn UI cho ảnh qua Redis và cập nhật kết quả vào PostgreSQL.

## 🧩 Các service/worker

- **label-platform-worker-gpt4**: Lấy job từ queue `label-platform-queue-gpt`, xử lý bằng model GPT-4 mock.
- **label-platform-worker-claude**: Lấy job từ queue `label-platform-queue-claude`, xử lý bằng model Claude mock.
- **label-platform-worker-gemini**: Lấy job từ queue `label-platform-queue-gemini`, xử lý bằng model Gemini mock.
- **label-platform-worker-result-consumer**: Lấy kết quả từ queue `label-platform-queue-result`, cập nhật vào bảng `image_samples` của PostgreSQL.

## Cấu trúc code

- `main.go`                : Worker chính (dùng chung cho các model, truyền biến môi trường để chọn model/queue)
- `types.go`               : Định nghĩa struct JobPayload, PredictionResultMessage, UIElement, Position
- `queue.go`               : Hàm thao tác Redis (BLPOP, RPUSH)
- `model_gpt.go`           : Mock model GPT-4
- `model_claude.go`        : Mock model Claude
- `model_gemini.go`        : Mock model Gemini
- `db.go`                  : Kết nối và cập nhật PostgreSQL (chỉ cho result-consumer)
- `Dockerfile`             : Build image cho từng worker

## Chạy thử bằng Docker Compose

1. Cấu hình biến môi trường trong file `.env` (xem ví dụ trong README này).
2. Build và chạy các container:

```bash
docker-compose up --build
```

## Biến môi trường

- `REDIS_ADDR`         : Địa chỉ Redis (vd: `redis:6379`)
- `REDIS_PASSWORD`     : Mật khẩu Redis (nếu có)
- `REDIS_DB`           : Số DB Redis (mặc định 0)
- `QUEUE_NAME`         : Tên queue Redis để lắng nghe (tùy worker)
- `RESULT_QUEUE`       : Tên queue kết quả (mặc định `label-platform-queue-result`)
- `MODEL_NAME`         : Tên model (gpt-4, claude, gemini)
- `POSTGRES_DSN`       : DSN kết nối PostgreSQL (chỉ cho result-consumer)

## Mô phỏng job push vào queue

```bash
redis-cli lpush label-platform-queue-gpt '{"image_id":"abc123","base64_image":"data:image/png;base64,..."}'
```

## Cập nhật DB

Worker result-consumer sẽ update trường `predicted_labels.<model>` trong bảng `image_samples`.

## Ví dụ push job và kiểm tra kết quả

1. Push job vào queue (ví dụ cho GPT-4):

```bash
redis-cli -h localhost lpush label-platform-queue-gpt '{"image_id":"abc123","base64_image":"data:image/png;base64,..."}'
```

2. Kiểm tra kết quả:

```bash
redis-cli -h localhost blpop label-platform-queue-result 0
```

## Lưu ý migrate DB

Bạn cần tạo bảng `image_samples` với trường `predicted_labels` kiểu JSONB trước khi chạy hệ thống.

```sql
CREATE TABLE image_samples (
  id uuid PRIMARY KEY,
  predicted_labels jsonb
);
```

## Đóng góp & mở rộng

- Dễ dàng thêm model mới bằng cách thêm file mock model và cấu hình biến môi trường.
- Code clean, dễ đọc, dễ mở rộng.
