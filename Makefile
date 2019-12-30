deezerdl:
	go build cmd/deezerdl.go

small:
	go build -ldflags="-w" cmd/deezerdl.go

clean:
	rm deezerdl
