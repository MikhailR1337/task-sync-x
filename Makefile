compile:
	go build -o build/main cmd/main.go
build-db:
	docker build -t task-sync-db ./services/postgres
run-db:
	docker run -p 5432:5432 -d task-sync-db
build-app:
	docker build -t task-sync-x .
run-app:
	docker run -p 3000:3000 -d task-sync-x
compose-up:
	docker-compose up -d
enter-db:
	docker exec -it task-sync-x-db-1 psql -U fanrik mydatabase
