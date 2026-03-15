PACKAGES_TESTS := $(shell go list ./...)
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
	go test -count=1 -race -v $(PACKAGES_TESTS)
	
coverage-terminal:
	go test -covermode=atomic -count=1 -race -coverprofile=coverage.out $(PACKAGES_TESTS)
	go tool cover -func=coverage.out

coverage-html:
	go test -count=1 -race -coverprofile=coverage.out $(PACKAGES_TESTS)
	go tool cover -html=coverage.out -o coverage.html

pprof-goroutine:
	go tool pprof -http=:6061 "http://localhost:6060/debug/pprof/goroutine"
pprof-cpu:
	go tool pprof -http=:6061 "http://localhost:6060/debug/pprof/profile?seconds=30"
pprof-heap:
	go tool pprof -http=:6061 "http://localhost:6060/debug/pprof/heap"