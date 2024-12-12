FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copiar arquivos de dependências
COPY go.mod go.sum ./

# Baixar dependências
RUN go mod download

# Copiar o código fonte
COPY . .

# Compilar a aplicação
RUN CGO_ENABLED=0 GOOS=linux go build -o dentist-api ./cmd/main.go

# Imagem final
FROM alpine:latest

WORKDIR /app

# Copiar o binário compilado da etapa anterior
COPY --from=builder /app/dentist-api .

# Expor a porta que a aplicação usa
EXPOSE 8080

# Comando para executar a aplicação
CMD ["./dentist-api"]