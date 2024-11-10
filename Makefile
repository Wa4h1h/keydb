.PHONY: build clean

local-deploy: build
	@cd build &&\
	export PORT="6000" &&\
    export LOG_LEVEL="debug" &&\
    export READ_TIMEOUT="30" &&\
	docker-compose up

build:
	GOOS=linux GOARCH=amd64 go build -o ./build/main ./cmd/dbserver/main.go

clean:
	rm ./build/main