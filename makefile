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
	go test -race ./...
	
coverage-terminal:
	go clean -testcache 
	go test -race -coverprofile=coverage.out ./...

coverage-html:
	make coverage-terminal
	go tool cover -html=coverage.out -o coverage.html

profiling:
	curl -o pprof/cpu.prof http://localhost:6061/debug/pprof/profile
	curl -o pprof/heap.prof http://localhost:6061/debug/pprof/heap
	curl -o pprof/goroutine.prof http://localhost:6061/debug/pprof/goroutine