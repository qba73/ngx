NGINX_PLUS_VERSION=r27
DOCKER_NETWORK?=test
DOCKER_NETWORK_ALIAS=nginx-plus-test
DOCKER_NGINX_PLUS?=nginx-plus
DOCKER_NGINX_PLUS_HELPER?=nginx-plus-helper

GO_CONTAINER=golang:1.19

export TEST_API_ENDPOINT=http://$(DOCKER_NGINX_PLUS):8080/api
export TEST_API_ENDPOINT_OF_HELPER=http://$(DOCKER_NGINX_PLUS_HELPER):8080/api
export TEST_UNAVAILABLE_STREAM_ADDRESS=$(DOCKER_NGINX_PLUS):8081

test: run-nginx-plus test-run configure-no-stream-block test-run-no-stream-block clean

lint:
	docker run --pull always --rm -v $(shell pwd):/ngx -w /ngx -v $(shell go env GOCACHE):/cache/go -e GOCACHE=/cache/go -e GOLANGCI_LINT_CACHE=/cache/go -v $(shell go env GOPATH)/pkg:/go/pkg golangci/golangci-lint:latest golangci-lint --color always run

build-debian:
	docker build --secret id=nginx-repo.crt,src=docker/nginx-repo.crt --secret id=nginx-repo.key,src=docker/nginx-repo.key --build-arg NGINX_PLUS_VERSION=$(NGINX_PLUS_VERSION) -t nginx-plus:$(NGINX_PLUS_VERSION) docker

build-alpine:
	docker build --secret id=nginx-crt,src=docker/nginx-crt --secret id=nginx-key,src=docker/nginx-key --build-arg NGINX_PLUS_VERSION=$(NGINX_PLUS_VERSION) -t nginx-plus:$(NGINX_PLUS_VERSION) -f docker/dockerfile.alpine .


run-nginx-plus:
	docker network create --driver bridge $(DOCKER_NETWORK)
	docker run --network=$(DOCKER_NETWORK) -d --name $(DOCKER_NGINX_PLUS) --network-alias=$(DOCKER_NETWORK_ALIAS) --rm -p 8080:8080 -p 8081:8081 nginx-plus:$(NGINX_PLUS_VERSION)
	docker run --network=$(DOCKER_NETWORK) -d --name $(DOCKER_NGINX_PLUS_HELPER) --network-alias=$(DOCKER_NETWORK_ALIAS) --rm -p 8090:8080 -p 8091:8081 nginx-plus:$(NGINX_PLUS_VERSION)

run-nginx:
	docker network create --driver bridge $(DOCKER_NETWORK)
	docker run --network=$(DOCKER_NETWORK) -d --name $(DOCKER_NGINX_PLUS) --network-alias=$(DOCKER_NETWORK_ALIAS) --rm -p 8080:8080 -p 8081:8081 nginx-plus:$(NGINX_PLUS_VERSION)
	docker run --network=$(DOCKER_NETWORK) -d --name $(DOCKER_NGINX_PLUS_HELPER) --network-alias=$(DOCKER_NETWORK_ALIAS) --rm -p 8090:8080 -p 8091:8081 nginx-plus:$(NGINX_PLUS_VERSION)



test-run:
	docker run --rm \
	--network=$(DOCKER_NETWORK) \
	-e TEST_API_ENDPOINT \
	-e TEST_API_ENDPOINT_OF_HELPER \
	-e TEST_UNAVAILABLE_STREAM_ADDRESS \
	-v $(shell pwd):/go/src/github.com/qba73/ngx \
	-w /go/src/github.com/qba73/ngx \
	$(GO_CONTAINER) /bin/sh -c "go test -count=1"

configure-no-stream-block:
	docker cp docker/nginx_no_stream.conf $(DOCKER_NGINX_PLUS):/etc/nginx/nginx.conf
	docker exec $(DOCKER_NGINX_PLUS) nginx -s reload

test-run-no-stream-block: configure-no-stream-block
	docker run --rm \
	--network=$(DOCKER_NETWORK) \
	-e TEST_API_ENDPOINT \
	-e TEST_API_ENDPOINT_OF_HELPER \
	-e TEST_UNAVAILABLE_STREAM_ADDRESS \
	-v $(shell pwd):/go/src/github.com/qba73/ngx \
	-w /go/src/github.com/qba73/ngx \
	$(GO_CONTAINER) /bin/sh -c "go test -count=1"

clean:
	-docker kill $(DOCKER_NGINX_PLUS)
	-docker kill $(DOCKER_NGINX_PLUS_HELPER)
	-docker network rm $(DOCKER_NETWORK)
