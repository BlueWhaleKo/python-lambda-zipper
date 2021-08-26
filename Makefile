build: build-linux build-windows

build-linux:
	goos=linux goarch=amd64 cgo_enabled=0 go build -installsuffix cgo -o  ./bin/python-packer ./src/main.go

build-windows:
	goos=windows goarch=386 cgo_enabled=0 go build -installsuffix cgo -o  ./bin/python-packer.exe ./src/main.go
