module riot-api

go 1.25.0

require (
	api-contracts v0.0.0-00010101000000-000000000000
	github.com/joho/godotenv v1.5.1
	google.golang.org/grpc v1.84.0
	google.golang.org/protobuf v1.36.12
)

require (
	github.com/golang/protobuf v1.5.4 // indirect
	golang.org/x/net v0.57.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
	golang.org/x/text v0.40.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260706201446-f0a921348800 // indirect
)

replace api-contracts => ../api-contracts
