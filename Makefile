run:
	go run src/main.go -p 5053 -d "127.0.0.1:5059"

build:
	go build -o bin/simple-choose-ad src/main.go
