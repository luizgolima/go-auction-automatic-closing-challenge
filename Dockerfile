FROM golang:1.26.3 as builder

WORKDIR /app

# Copia arquivos de dependências
COPY go.mod go.sum ./
RUN go mod download

# Copia o restante do código
COPY . .

# Build do binário estático
RUN GOOS=linux CGO_ENABLED=0 go build -o auction cmd/auction/main.go

# Imagem final mínima
FROM scratch

# Copia o binário
COPY --from=builder /app/auction /auction

# MANTÉM A ESTRUTURA DE PASTAS QUE O CÓDIGO ESPERA PARA O .ENV
COPY --from=builder /app/cmd/auction/.env /cmd/auction/.env

ENTRYPOINT ["./auction"]
