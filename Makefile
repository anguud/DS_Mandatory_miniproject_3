gen:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/service.proto

cli1: 
	go run client/client.go

cli2: 
	go run client/client.go

cli3: 
	go run client/client.go

repl1: 
	go run server/server.go -port 9080

repl2: 
	go run server/server.go -port 9081

repl3: 
	go run server/server.go -port 9082
