BUILD_ROOT = ./bin

build:
	go build -o $(BUILD_ROOT)/notadash \
		-ldflags "-X main.VERSION $(shell cat notadash/VERSION)" \
		notadash/*.go

build-mon:
	go build -o $(BUILD_ROOT)/notadash-mon \
		-ldflags "-X main.VERSION $(shell cat notadash-mon/VERSION)" \
		notadash-mon/*.go

test-deps:
	go get github.com/stretchr/testify
	go get golang.org/x/tools/cmd/cover

build-deps:
	# TODO (boldfield) :: Generalize this
	go get github.com/codegangsta/cli
	go get github.com/ryanuber/columnize
	go get code.google.com/p/gcfg
	go get github.com/fsouza/go-dockerclient
	go get github.com/gambol99/go-marathon
	go get github.com/boldfield/go-mesos
	go get golang.org/x/crypto/ssh/terminal
