FROM golang:alpine AS backend-builder

WORKDIR /app

COPY . .

RUN go mod download

RUN GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o main ./cmd


FROM alpine:latest AS runner

WORKDIR /app

# 安装 gcc 和其他必要的构建工具
RUN apk add --no-cache gcc

COPY --from=backend-builder /app/main .

CMD ["./main"]