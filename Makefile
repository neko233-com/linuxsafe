.PHONY: build test test-docker clean install docs

build:
	go build -o linuxsafe .

test:
	go test ./... -v -count=1

test-docker: test-docker-ubuntu test-docker-debian test-docker-centos test-docker-alpine

test-docker-ubuntu:
	@echo "Testing on Ubuntu..."
	docker build -f test/docker/Dockerfile.ubuntu -t linuxsafe-test-ubuntu .
	docker run --rm linuxsafe-test-ubuntu

test-docker-debian:
	@echo "Testing on Debian..."
	docker build -f test/docker/Dockerfile.debian -t linuxsafe-test-debian .
	docker run --rm linuxsafe-test-debian

test-docker-centos:
	@echo "Testing on CentOS..."
	docker build -f test/docker/Dockerfile.centos -t linuxsafe-test-centos .
	docker run --rm linuxsafe-test-centos

test-docker-alpine:
	@echo "Testing on Alpine..."
	docker build -f test/docker/Dockerfile.alpine -t linuxsafe-test-alpine .
	docker run --rm linuxsafe-test-alpine

install:
	go install .

clean:
	rm -f linuxsafe linuxsafe.exe
	docker rmi linuxsafe-test-ubuntu linuxsafe-test-debian linuxsafe-test-centos linuxsafe-test-alpine 2>/dev/null || true

docs:
	@echo "Generating docs..."
	@mkdir -p docs/site
	@cp -r docs/static/* docs/site/ 2>/dev/null || true

release-snapshot:
	goreleaser release --snapshot --clean
