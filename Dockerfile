# Sử dụng Go release candidate trên Alpine (cho hiệu suất tốt hơn)
FROM golang:rc-alpine AS builder

# Đặt thư mục làm việc
WORKDIR /app

# Copy go.mod và go.sum để tận dụng cache
COPY go.mod go.sum ./

# Tải dependencies
RUN go mod download

# Copy toàn bộ mã nguồn vào container
COPY . .

# Build ứng dụng (tạo file thực thi)
RUN go build -o main .

# Tạo container chạy ứng dụng (multi-stage build để giảm kích thước)
FROM alpine:latest  

WORKDIR /root/

# Copy file thực thi từ builder stage
COPY --from=builder /app/main .

# Xuất port (Render sẽ tự động nhận diện)
EXPOSE 8080

# Chạy ứng dụng
CMD ["./main"]