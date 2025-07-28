# --- Builder Stage ---
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /loadtester ./cmd

# --- Final Stage ---
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /loadtester .
ENTRYPOINT ["./loadtester"]