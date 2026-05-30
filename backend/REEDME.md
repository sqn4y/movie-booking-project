# Backend Service

REST-сервис для операций с фильмами и бронированиями.

## Запуск

Команды нужно выполнять из директории `backend`.

```bash
docker compose up --build
go run ./cmd/app
```

Сервис будет доступен на `http://localhost:8080`.

Swagger UI доступен на `http://localhost:8080/swagger/index.html`.

## Тесты

```bash
make test
```

Дополнительные команды:

```bash
make test100
make race
make cover
```

## Генерация моков

```bash
make gen
```

Команда генерирует моки для интерфейсов repository и service:

- `internal/repository/mock`
- `internal/service/mock`

## Ручки

- `POST /api/v1/movie` - создать фильм
- `GET /api/v1/movies` - получить список фильмов
- `POST /api/v1/booking` - создать бронь
- `DELETE /api/v1/booking/{id}` - удалить бронь
- `PUT /api/v1/booking/{id}` - изменить бронь

## Пример создания фильма

```bash
curl -X POST http://localhost:8080/api/v1/movie \
  -H 'Content-Type: application/json' \
  -d '{
      "title": "Оппенгеймер",
      "director": "Кристофер Нолан",
      "duration": 10800,
      "description": "История американского физика Роберта Оппенгеймера",
      "genre_ids": [1,3],
      "age_rating": 18,
      "release_date": "2024-12-01T00:00:00Z"
  }'
```

## Пример создания брони

```bash
curl -X POST http://localhost:8080/api/v1/booking \
  -H 'Content-Type: application/json' \
  -d '{
      "movie_id": 1,
      "seats": ["A1", "A2"]
  }'
```
