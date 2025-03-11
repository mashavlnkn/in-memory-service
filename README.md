# In-Memory Task Service

## 📌 Описание проекта
Это простой сервис управления задачами, реализованный на Go с использованием фреймворка **Fiber** и логгера **Zap**.  
Сервис использует **in-memory (хранение в оперативной памяти)** для хранения данных, что делает его быстрым, но без постоянного сохранения информации.

## 🚀 Функциональность
- Создание задачи (`POST /tasks`)
- Получение списка всех задач (`GET /tasks`)
- Получение задачи по ID (`GET /tasks/:id`)
- Обновление задачи (`PUT /tasks/:id`)
- Удаление задачи (`DELETE /tasks/:id`)

## 🛠️ Технологии
- **Go** (1.21+)
- **Fiber** (веб-фреймворк)
- **Zap** (логирование)
- **Validator** (валидация входных данных)

## 📦 Установка и запуск

### 1️⃣ Клонируем репозиторий
```sh
git clone https://github.com/mashavlnkn/in-memory-service.git
cd in-memory-service
```
### 2️⃣ Устанавливаем зависимости
```js
go mod tidy
```
### 3️⃣ Запускаем сервис
```js
go run cmd/main.go

```
##  Примеры запросов
### 1. Создать задачу
```js
curl -X POST http://localhost:8080/tasks \
     -H "Content-Type: application/json" 
     -d '{
           "title": "Сделать задание",
           "description": "Написать сервис на Go",
           "status": "pending"
         }'

```
### 2. Получить все задачи
```js
curl -X GET http://localhost:8080/tasks

```
### 3. Получить задачу по ID
```js
curl -X GET http://localhost:8080/tasks/1
```
### 4. Обновить задачу
```js
curl -X PUT http://localhost:8080/tasks/1 \
-H "Content-Type: application/json" \
-d '{
"title": "Обновленное задание",
"description": "Теперь с новыми требованиями",
"status": "in_progress"
}'
```

### 5. Удалить задачу
```
curl -X DELETE http://localhost:8080/tasks/1
```

### Дополнительно
Данный сервис не использует базу данных и хранит данные в оперативной памяти.
При перезапуске все данные будут потеряны.
Можно доработать сервис, добавив поддержку SQLite, PostgreSQL или MongoDB.
### Контакты
 #### Автор: @mashavlnkn
#### GitHub: https://github.com/mashavlnkn