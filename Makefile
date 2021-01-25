.PHONY: linux_64 freebsd_64 openbsd_64 darwin_64 clean

all: linux freebsd openbsd darwin

linux:
	# Building Linux64
	@ CGO_ENABLED=0 GOOS=linux   GOARCH=amd64 go build -o build/statictls_linux64    cmd/statictls/main.go

freebsd:
	# Building FreeBSD64
	@ CGO_ENABLED=0 GOOS=freebsd GOARCH=amd64 go build -o build/statictls_freebsd_64 cmd/statictls/main.go

openbsd:
	# Building OpenBSD64
	@ CGO_ENABLED=0 GOOS=openbsd GOARCH=amd64 go build -o build/statictls_openbsd64  cmd/statictls/main.go

darwin:
	# Building Darwin64
	@ CGO_ENABLED=0 GOOS=darwin  GOARCH=amd64 go build -o build/statictls_darwin64   cmd/statictls/main.go

clean:
	@rm ./build/*
