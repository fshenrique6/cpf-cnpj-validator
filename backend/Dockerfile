# Etapa de build
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o cpf-cnpj-validator .

# Etapa final (imagem enxuta, só com o binário)
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/cpf-cnpj-validator .

EXPOSE 8080

CMD ["./cpf-cnpj-validator"]