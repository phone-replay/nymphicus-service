# Etapa 1: Compilação
FROM golang:1.20 as builder

# Definindo o diretório de trabalho
WORKDIR /app

# Copiando os arquivos go.mod e go.sum
COPY go.mod go.sum ./

# Instalando as dependências
RUN go mod download

# Copiando o restante dos arquivos da aplicação
COPY . .

# Compilando a aplicação
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

# Etapa 2: Execução
FROM alpine:latest

# Instalando o ca-certificates para SSL
RUN apk --no-cache add ca-certificates

# Definindo o diretório de trabalho
WORKDIR /root/

# Copiando o binário compilado e o arquivo de configuração da etapa de build
COPY --from=builder /app/main .
COPY --from=builder /app/config/config-local.yml ./config/config-local.yml

# Expondo a porta que a aplicação irá rodar
EXPOSE 8080

# Comando para iniciar a aplicação
CMD ["./main"]
