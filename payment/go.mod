module github.com/space-wanderer/microservices/payment

replace github.com/space-wanderer/microservices/shared => ../shared

go 1.24.4

require (
	github.com/google/uuid v1.6.0
	github.com/space-wanderer/microservices/shared v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.73.0
)

require (
	golang.org/x/net v0.40.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.25.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250324211829-b45e905df463 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
)
