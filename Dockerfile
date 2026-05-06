FROM golang:1.26-alpine AS builder
WORKDIR /app

# 1. Download dependencies first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# 2. Copy only the source code directories needed for the build
COPY cmd/ cmd/
COPY internal/ internal/

# 3. Build the static binary with optimizations
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o api cmd/server/main.go

FROM scratch
WORKDIR /app

# 4. Copy the binary from the builder stage
COPY --from=builder /app/api /app/api

# 5. Copy ONLY the production-ready resource files explicitly
# This prevents including huge source datasets (like references.json.gz) in the final image
COPY resources/ivf.bin /app/resources/ivf.bin
COPY resources/mcc_risk.json /app/resources/mcc_risk.json
COPY resources/normalization.json /app/resources/normalization.json

CMD ["/app/api"]
