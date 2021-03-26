# Собрать образ приложения из Dokcerfile : docker image build --target prod -t jwtauthv2_image . 
# Запуск приложения и всех сервисов в докере : docker-compose up --remove-orphans --build jwtauthv2
# Запуск бд : docker-compose up -d db adminer && make migrateup


export DSN_DB=postgresql://postgres:example@localhost:5432/jwtauthv2?sslmode=disable
export SECRET=SECRET
export PORT=:8000

# Локальный запуск 
start: 
	DSN_DB=$$DSN_DB SECRET=$$SECRET PORT=$$PORT go run ./cmd/main.go


migrateup:
	migrate -path ./migration -database "postgresql://postgres:example@localhost:5432/jwtauthv2?sslmode=disable" -verbose up
	

migratedown: 
	migrate -path ./migration -database "postgresql://postgres:example@localhost:5432/jwtauthv2?sslmode=disable" -verbose down

test.e2e:
	DSN_DB=$$DSN_DB SECRET=$$SECRET PORT=$$PORT go test -v ./test/e2e/