FROM golang:1.23.4 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .

WORKDIR /app/api/cmd
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main .

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/main .

EXPOSE 8000

CMD ["./main"]
