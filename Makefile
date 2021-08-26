build:
	goos=linux goarch=amd64 cgo_enabled=0 go build -installsuffix cgo -o  ./bin/python-packer ./src/main.go
