FROM golang:alpine AS builder

WORKDIR /app

# 复制 go.mod 和 go.sum 以利用缓存
COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o main ./cmd


FROM alpine:3.21 AS runner

WORKDIR /app

# 安装 gcc musl-dev(C语言必要工具) python3
RUN apk add --no-cache gcc musl-dev python3 g++

COPY --from=builder /app/main .
COPY ./data/* ./data/
EXPOSE 3000

CMD ["./main"]
