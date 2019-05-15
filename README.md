[![Build Status](https://travis-ci.org/socrata-platform/notadash.svg?branch=master)](https://travis-ci.org/socrata-platform/notadash)

# notadash

notadash is a dashboard/monitoring service interface for the Mesos/Marathon/Docker stack.  It provides a CLI dashboard
for monitoring the current state of the stack as well as provide a simple tool to use in conjunction with monitoring services
such as [Sensu](http://sensuapp.org/).

# Building

1. Install go
1. Check out the repo into `$GOPATH/src/github.com/socrata-platform/notadash`  
don't have a `$GOPATH` set it defaults to `~/go`
1. run `make build-deps test-deps test build build-mon`
1. If successful, you should see `notadash` and `notadash-mon` in the `bin` subdirectory

## notadash

Docs TBD

## notadash-mon

Docs TBD
