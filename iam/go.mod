module github.com/space-wanderer/microservices/iam

go 1.24.4

require (
	github.com/caarlos0/env/v11 v11.3.1
	github.com/jackc/pgx/v5 v5.7.5
	github.com/joho/godotenv v1.5.1
	google.golang.org/protobuf v1.36.8
)

require (
	github.com/gomodule/redigo v1.9.2
	github.com/space-wanderer/microservices/platform v0.0.0
	github.com/space-wanderer/microservices/shared v0.0.0
	go.uber.org/zap v1.27.0
	google.golang.org/grpc v1.75.0
)

require (
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/crypto v0.41.0 // indirect
	golang.org/x/net v0.43.0 // indirect
	golang.org/x/sync v0.16.0 // indirect
	golang.org/x/sys v0.35.0 // indirect
	golang.org/x/text v0.28.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250825161204-c5933d9347a5 // indirect
)

replace (
	github.com/space-wanderer/microservices/platform => ../platform
	github.com/space-wanderer/microservices/shared => ../shared
)
