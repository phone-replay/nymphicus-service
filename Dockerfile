FROM golang:1.22.2 as builder
WORKDIR /nymphicus-service
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main .
FROM debian:bullseye-slim
WORKDIR /root/
COPY --from=builder /nymphicus-service/main .
CMD ["./main"]