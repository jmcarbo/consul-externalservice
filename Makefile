test:
	go test -v .

build: cli/consul-externalservice.go
	gox -output "bin/consul-externalservice_{{.OS}}_{{.Arch}}" -os "linux darwin" -arch "amd64" ./cli

startconsul:
	rm -rf /tmp/consul
	consul agent -server -bootstrap -data-dir /tmp/consul

stopconsul:
	killall -TERM consul


build_docker: cli/consul-externalservice.go
	docker build -t jmcarbo/consul-externalservice .

run_docker: cli/consul-externalservice.go
	docker run -ti --rm --name node1 jmcarbo/consul-externalservice /bin/bash

