# Этап сборки
FROM node:lts AS build

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем package.json и package-lock.json (или yarn.lock)
COPY package*.json ./

# Устанавливаем зависимости
RUN npm install

# Копируем все исходные файлы приложения
COPY . .

# Строим приложение для продакшн
RUN npm run build

# Этап для продакшн-сервера
FROM node:lts AS production

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем только билд из первого этапа
COPY --from=build /app/dist /app/dist

# Устанавливаем serve для запуска продакшн-сервера
RUN npm install -g serve

# Открываем порт 3000
EXPOSE 3000

# Запускаем приложение с помощью serve
CMD ["serve", "-s", "dist", "-l", "3000"]
