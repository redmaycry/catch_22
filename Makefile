# port for main server
port := 5053
moc_server_address := 127.0.0.1:5059

run:
	go run cmd/main.go -p $(port) -d "$(moc_server_address)"

test-ip:
	go run cmd/main.go -p $(port) -d "$(moc_server_address),localhost:5059"

test-port:
	go run cmd/main.go -p $(port) -d "$(moc_server_address),127.0.0.1:as"

test-port-max:
	go run cmd/main.go -p $(port) -d "$(moc_server_address),127.0.0.1:65537"

test-port-endpoint:
	go run cmd/main.go -p $(port) -d "127.0.0.1:9001/bid_request"

build:
	go build -o bin/simple-choose-ad cmd/main.go

start-moc-server:
	@echo "[!] Starting moc server on 127.0.0.1:5059..."
	@nohup go run internal/moc_server.go> /dev/null 2>&1 & echo $! > run.pid

stop-moc-server:
	@echo "[!] Killing moc server"
	@cat run.pid | xargs kill -9
	@rm run.pid

test-server:
	@echo "Testing server..."
	@cd "cmd/client_server/"; \
	go test

tests: start-moc-server
	@$(MAKE) test-server
	@$(MAKE) stop-moc-server
