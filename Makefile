AUBIO_PREFIX := $(shell brew --prefix aubio)
CGO_CFLAGS := -I$(AUBIO_PREFIX)/include
CGO_LDFLAGS := -L$(AUBIO_PREFIX)/lib

build:
	CGO_CFLAGS="$(CGO_CFLAGS)" CGO_LDFLAGS="$(CGO_LDFLAGS)" go build ./...

run:
	CGO_CFLAGS="$(CGO_CFLAGS)" CGO_LDFLAGS="$(CGO_LDFLAGS)" go run .

clean:
	go clean ./...
