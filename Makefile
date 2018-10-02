.PHONY: shell build test clean docker-up
SHELL=/bin/bash -o pipefail

BIN=pagerbot
GO_FLAGS=-ldflags "-extldflags '-static'"

DOCKER_FLAGS=-u $$(id -u):$$(id -g)
RUNNER=docker-compose exec -T $(DOCKER_FLAGS) runner

shell: docker-up
	docker-compose exec $(DOCKER_FLAGS) runner bash

test: docker-up
	@$(RUNNER) go test -v ./...

clean:
	@docker-compose down --rmi all

###

docker-up:
	@(env -i bash --noprofile --norc -c '. platform/secrets/ci.env; env') | grep -v '^PWD=' > .ci-runner.env
	@env | ( grep DOCKER_ || true ) >> .ci-runner.env
	@FIXUID=$$(id -u) FIXGID=$$(id -g) docker-compose up -d
	@$(RUNNER) go get github.com/mitchellh/gox

build: docker-up
	@$(RUNNER) gox -osarch "linux/amd64" $(GO_FLAGS) -output "dist/{{.OS}}_{{.Arch}}/$(BIN)"
