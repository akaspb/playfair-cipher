EXEC_PATH=bin/playfair

build:
	go build -o ${EXEC_PATH} cmd/cipher-test/main.go

run:
	./${EXEC_PATH}

full: build run
