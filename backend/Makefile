.PHONY: test test100 race cover gen

MOCKGEN := go run go.uber.org/mock/mockgen

test:
	go test -v -count=1 ./...

test100:
	go test -v -count=100 ./...

race:
	go test -v -race -count=1 ./...

cover:
	go test -short -count=1 -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
	rm -f coverage.out

gen:
	mkdir -p internal/repository/mock internal/service/mock
	$(MOCKGEN) -source=internal/repository/movie.go -destination=internal/repository/mock/mock_movie_repository.go -package mock
	$(MOCKGEN) -source=internal/repository/booking.go -destination=internal/repository/mock/mock_booking_repository.go -package mock
	$(MOCKGEN) -source=internal/service/movie.go -destination=internal/service/mock/mock_movie_service.go -package mock
	$(MOCKGEN) -source=internal/service/booking.go -destination=internal/service/mock/mock_booking_service.go -package mock
