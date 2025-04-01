# FROM golang:1.24

# WORKDIR /app

# COPY . .

# #RUN go mod tidy && go mod download

# # Установка air с использованием нового пути
# RUN go install github.com/air-verse/air@latest && \
#     ls -l $(go env GOPATH)/bin && \
#     cp $(go env GOPATH)/bin/air /usr/local/bin/air && \
#     ls -l /usr/local/bin && \
#     chmod +x /usr/local/bin/air

# EXPOSE 8080

# CMD ["sh", "-c", "$(go env GOPATH)/bin/air"]

# Устанавливаем Go и необходимые инструменты
FROM golang:1.24

# Установка рабочей директории внутри контейнера
WORKDIR /app

# Копирование исходного кода в контейнер
COPY . .

# Обновление зависимостей Go
RUN go mod tidy
RUN go mod download

# Установка Air для автоматической перезагрузки
RUN go install github.com/air-verse/air@latest

# Установка Delve для отладки
RUN go install github.com/go-delve/delve/cmd/dlv@latest

# Открытие портов для работы сервера и отладчика
EXPOSE 8080 2345