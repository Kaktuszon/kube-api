APP_NAME = kube-api

build:
	go build -o bin/${APP_NAME}

run: build
	./bin/${APP_NAME} -path "bin/config.yaml"

archive: build
	tar -cvf bin/${APP_NAME}.tar bin/${APP_NAME}