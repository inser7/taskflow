FROM node:lts-alpine

# Переключаемся в root для установки зависимостей
USER root

WORKDIR /app

# # Копируем package.json перед установкой зависимостей
# COPY package.json package-lock.json ./

# # Даем права пользователю node
# RUN chown -R node:node /app

# # Переключаемся обратно на node
# USER node

# # Устанавливаем зависимости
# RUN npm install

# # Копируем оставшиеся файлы
# COPY --chown=node:node . .

# # Запуск
# CMD ["npm", "run", "dev"]