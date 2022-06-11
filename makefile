.DEFAULT_GOAL := build

fmt:
	go fmt ./...
.PHONY:fmt

lint: fmt
	golint ./...
.PHONY:lint

vet: fmt
	go vet ./...
.PHONY:vet

build: vet
	go build -o bin/tankerkoenig-to-influxdb ./...
.PHONY:build

build_arm_7:
	GOOS=linux GOARCH=arm GOARM=7 go build -o bin/tankerkoenig-to-influxdb-linux-arm-7 ./...
.PHONY:build_arm_7
