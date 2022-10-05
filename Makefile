# port for main server
port := 5053
moc_server_address := 127.0.0.1:5059

bold := \033[1m
normal := \033[0m
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
	@echo "[!] Starting moc server on $(moc_server_address) ..."
	@go run internal/moc_server.go -l $(moc_server_address) &


stop-moc-server:
	@echo "[!] Killing moc server"
	@curl -s -o /dev/null "$(moc_server_address)/exit" &

test-server:
	@echo "Testing server..."
	@$(MAKE) start-moc-server
	@cd "cmd/client_server/"; \
	go test
	@$(MAKE) stop-moc-server

tests:
	@$(MAKE) test-server

help:
	@echo "$(bold)Makefile commands$(normal)"
	@echo "-----------------"
	@echo "$(bold)make build$(normal)   : will build the project"
	@echo "$(bold)make tests$(normal)   : run tests for the project"
	@echo "$(bold)make run$(normal)     : will run the project"
	@echo ""
	@echo "$(bold)OS commands$(normal)"
	@echo "-----------"
	@echo "start server at PORT with 'IP:PORT' list of partners:"
	@echo "$(bold)./bin/simple-choose-ad -p PORT -d 'IP:PORT'$(normal)"
