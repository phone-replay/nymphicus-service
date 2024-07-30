FROM golang:1.20 as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
COPY --from=builder /app/config/config-local.yml ./config/config-local.yml
COPY --from=builder /app/config/config-production.yml ./config/config-production.yml
EXPOSE 8080
CMD ["./main"]
