
```bash
docker compose up --build
go run ./cmd/app
```

Команды нужно выполнять из директории `backend`.

Сервис будет доступен на `http://localhost:8080`.

## Ручки

- `POST /api/v1/movie` - создать фильм
- `GET /api/v1/movies` - получить список фильмов
- `POST /api/v1/booking` - создать бронь
- `DELETE /api/v1/booking/{id}` - удалить бронь
- `PUT /api/v1/booking/{id}` - изменить бронь

## Пример создания

```bash

      "age_rating": 18,
      "release_date": "2024-12-01T00:00:00Z"
  }'
```
