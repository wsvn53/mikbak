darwin:
	go build -o mikbak-darwin

linux:
	GOOS=linux GOARCH=amd64 go build -o mikbak-linux
