.buildall: all
all:build 


generate:
	@go generate ./...
	@echo "[OK] Files added to embed box!"

build: generate
	CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o goinventory ./cmd/app/*.go
	@echo "[OK] App binary was created!"

run:
	@./build/server
