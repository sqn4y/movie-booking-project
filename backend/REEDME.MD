# Backend Service

REST-сервис для CRUD-операций бронирования фильмов.

## Запуск

```bash
docker compose up --build
```
в директории /backend/cmd/app выполнить команду `go run main.go`

Сервис будет доступен на `http://localhost:8080`.


## Ручки

- `POST /api/v1/movie` - создать фильм
- `GET /api/v1/movies` - получить список фильмов
- 
- `POST /api/v1/booking` - создать бронь
- `DELETE /api/v1/booking/{id}` - удалить бронь
- `PUT /api/v1/booking/{id}` - изменить бронь


## Пример создания

```bash
curl -X POST http://localhost:8080/api/v1/movie \
  -H 'Content-Type: application/json' \
  -d '{
      "title": "Оппенгеймер",
      "director": "Кристофер Нолан",
      "duration": 10800,
      "description": "История американского физика Роберта Оппенгеймера, который руководил Манхэттенским проектом",
      "genre_ids": [1,3],
      "age_rating": 18,
      "release_date": "2024-12-01T00:00:00Z"
  }'
```