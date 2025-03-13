# Sử dụng Go 1.23
FROM golang:1.23 AS builder

# Đặt thư mục làm việc
WORKDIR /app

# Copy go.mod và go.sum để tận dụng cache
COPY go.mod go.sum* ./

# Tải dependencies
RUN go mod download

# Copy toàn bộ mã nguồn vào container
COPY . .

# Build ứng dụng (tạo file thực thi)
RUN go build -o main .

# Tạo container chạy ứng dụng (multi-stage build để giảm kích thước)
FROM debian:bookworm-slim

WORKDIR /root/

# Cài đặt các gói cần thiết
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

# Tạo thư mục credentials
RUN mkdir -p /root/credentials

# Copy file thực thi từ builder stage
COPY --from=builder /app/main .

# Copy thư mục credentials (cần cho Google Sheets API)
COPY --from=builder /app/credentials/ /root/credentials/

# Xuất port (Render sẽ tự động nhận diện)
EXPOSE 8000

# Chạy ứng dụng
CMD ["/root/main"]