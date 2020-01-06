deezerdl:
	go build cmd/deezerdl.go

small:
	go build -ldflags="-w" cmd/deezerdl.go

releases:
	mkdir releases
	GOOS="linux" GOARCH="amd64" go build -o "releases/deezerdl-linux-amd64" cmd/deezerdl.go
	GOOS="linux" GOARCH="386" go build -o "releases/deezerdl-linux-i386" cmd/deezerdl.go
	GOOS="windows" GOARCH="amd64" go build -o "releases/deezerdl-windows-amd64.exe" cmd/deezerdl.go
	GOOS="windows" GOARCH="386" go build -o "releases/deezerdl-windows-i386.exe" cmd/deezerdl.go

clean:
	rm -f deezerdl
	rm -rf releases
