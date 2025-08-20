module github.com/space-wanderer/microservices/notification

go 1.24.4

require (
	github.com/caarlos0/env/v11 v11.2.2
	github.com/joho/godotenv v1.5.1
	github.com/space-wanderer/microservices/platform v0.0.0
	github.com/space-wanderer/microservices/shared v0.0.0
)

replace (
	github.com/space-wanderer/microservices/platform => ../platform
	github.com/space-wanderer/microservices/shared => ../shared
)
