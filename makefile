docker-monitor:
	docker run --rm -it -v /var/run/docker.sock:/var/run/docker.sock -v ./config:/.config/jesseduffield/lazydocker lazyteam/lazydocker

docker-inspect-leak
	go tool pprof -http :9999 -edgefraction 0 -nodefraction 0 -nodecount 100000