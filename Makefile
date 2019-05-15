BUILD_ROOT = ./bin

test:
	go test ./... -v -cover

build:
	go build -o $(BUILD_ROOT)/notadash \
		-ldflags "-X main.VERSION=$(shell cat notadash/VERSION)" \
		notadash/*.go

build-mon:
	go build -o $(BUILD_ROOT)/notadash-mon \
		-ldflags "-X main.VERSION=$(shell cat notadash-mon/VERSION)" \
		notadash-mon/*.go

test-deps:
	go get github.com/stretchr/testify
	go get golang.org/x/tools/cmd/cover

build-deps:
	# TODO (boldfield) :: Generalize this
	go get github.com/marpaia/graphite-golang
	go get github.com/codegangsta/cli
	go get github.com/ryanuber/columnize
	go get github.com/scalingdata/gcfg
	go get github.com/fsouza/go-dockerclient
	go get github.com/gambol99/go-marathon
	go get github.com/socrata-platform/go-mesos
	go get golang.org/x/crypto/ssh/terminal
	go get github.com/behance/go-chronos/chronos
