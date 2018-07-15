
NAME := api

default: start

.PHONY: setup
setup:
	go get -u github.com/golang/dep/cmd/dep
	$(MAKE) dep

.PHONY: dep
dep:
	@dep ensure -v

.PHONY: start
start:
	mkdir -p build
	go build -o build/${NAME} github.com/sandyleo26/sydney_weather
	./build/${NAME}

.PHONY: clean
clean:
	rm -rf build/