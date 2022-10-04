run:
	go run cmd/main.go -p 5053 -d "127.0.0.1:5059"

test-ip:
	go run cmd/main.go -p 5053 -d "127.0.0.1:5059,localhost:5059"

test-port:
	go run cmd/main.go -p 5053 -d "127.0.0.1:5059,127.0.0.1:as"

test-port-max:
	go run cmd/main.go -p 5053 -d "127.0.0.1:5059,127.0.0.1:65537"

test-port-endpoint:
	go run cmd/main.go -p 5053 -d "127.0.0.1:9001/bid_request"

build:
	go build -o bin/simple-choose-ad cmd/main.go

tests:
