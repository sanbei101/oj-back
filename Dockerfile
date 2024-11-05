FROM golang:alpine AS backend-builder

WORKDIR /app

COPY . .

RUN go mod download

RUN GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o main ./cmd


FROM alpine:latest AS runner

WORKDIR /app

COPY --from=backend-builder /app/main .
COPY --from=frontend-builder /app/pocketbase-vue/dist ./dist

CMD ["./main"]