FROM golang:alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o main ./cmd


FROM alpine:3.21 AS runner

WORKDIR /app

# 安装 gcc 和其他必要的构建工具
RUN apk add --no-cache gcc musl-dev python3

COPY --from=builder /app/main .

EXPOSE 3000

CMD ["./main"]
