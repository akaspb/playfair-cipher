BIN_PATH=bin
EXEC_PATH=bin/playfair
EXEC_PATH_WIN=bin/playfair.exe

build:
	mkdir -p ${BIN_PATH}
	go build -o ${EXEC_PATH} cmd/main.go

run:
	./${EXEC_PATH}

full: build run

conf:
	go run cmd/make-config/en/main.go

conf-ru:
	go run cmd/make-config/ru/main.go

build-win:
	mkdir -p ${BIN_PATH}
	GOOS=windows GOARCH=386 go build -o ${EXEC_PATH_WIN} cmd/main.go
