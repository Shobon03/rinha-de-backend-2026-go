FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o api cmd/server/main.go
FROM scratch
WORKDIR /app
COPY --from=builder /app/api /app/api
COPY resources/ /app/resources/
CMD ["/app/api"]
