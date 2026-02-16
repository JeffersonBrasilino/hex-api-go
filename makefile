TIMESTAMP := $(shell date +'%Y-%m-%d_%H-%M-%S')
docker-monitor:
	docker run --rm -it -v /var/run/docker.sock:/var/run/docker.sock -v ./config:/.config/jesseduffield/lazydocker lazyteam/lazydocker

deps:
	go mod vendor
	go clean -modcache

start-dev:
	echo "Verificando dependências..."
	@if [ ! -d vendor ]; then \
		echo "Diretório 'vendor' não encontrado. Executando 'make deps'..."; \
		make deps; \
	fi
	docker compose up -d

test:
	go clean -testcache
	go test -race -v ./...
	
coverage-terminal:
	go clean -testcache 
	go test -race -cover ./...

coverage-html:
	go test -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

pprof-goroutine:
	go tool pprof -http=:6061 "http://localhost:6060/debug/pprof/goroutine"
pprof-cpu:
	go tool pprof -http=:6061 "http://localhost:6060/debug/pprof/profile?seconds=30"
pprof-heap:
	go tool pprof -http=:6061 "http://localhost:6060/debug/pprof/heap"