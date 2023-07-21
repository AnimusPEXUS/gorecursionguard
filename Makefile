export GONOPROXY=github.com/AnimusPEXUS/*

all: get

get:
		$(MAKE) -C examples/01 get
		go get -u -v "./..."
		go mod tidy

build:
		$(MAKE) -C examples/01 build
		go build
