FROM golang:1.21 AS builder

WORKDIR /app

COPY . .

RUN go mod tidy && go build -o taskflow

FROM debian:bullseye
WORKDIR /root/
COPY --from=builder /app/taskflow .
EXPOSE 8080
CMD ["./taskflow"]
