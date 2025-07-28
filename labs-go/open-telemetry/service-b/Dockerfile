# Estágio de build
FROM golang:1.23-alpine as builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /main .

# Estágio final (imagem de produção)
FROM alpine:latest
WORKDIR /app
COPY --from=builder /main .
# O Serviço B roda na porta 8081
EXPOSE 8081
ENTRYPOINT ["/app/main"]
