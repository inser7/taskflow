FROM golang:1.24

WORKDIR /app

COPY . .

#RUN go mod tidy && go mod download

# Установка air с использованием нового пути
RUN go install github.com/air-verse/air@latest && \
    ls -l $(go env GOPATH)/bin && \
    cp $(go env GOPATH)/bin/air /usr/local/bin/air && \
    ls -l /usr/local/bin && \
    chmod +x /usr/local/bin/air

EXPOSE 8080

CMD ["sh", "-c", "$(go env GOPATH)/bin/air"]