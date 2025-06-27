APP_EXECUTABLE=bin/goputer

build:
	@GOARCH=amd64 GOOS=linux go build -o ${APP_EXECUTABLE}-linux.exe main.go
	@GOARCH=amd64 GOOS=windows go build -o ${APP_EXECUTABLE}-windows.exe main.go
	@GOARCH=arm64 GOOS=darwin go build -o ${APP_EXECUTABLE}-arm64.exe main.go

run:
	@go run main.go

clean:
	@go clean
	@rm ${APP_EXECUTABLE}-linux.exe
	@rm ${APP_EXECUTABLE}-windows.exe
	@rm ${APP_EXECUTABLE}-arm64.exe
