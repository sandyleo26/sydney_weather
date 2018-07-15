Sydney Weather
--

## Setup
You need to have `go` installed and `GOPATH` set correctly.

```
go get -u github.com/sandyleo26/sydney_weather

cd $GOPATH/src/github.com/sandyleo26/sydney_weather

make setup

make
```

## Usage
I use `curl` for testing. But tools like postman also works.

```
curl http://localhost:8080/v1/weather?city=sydney

```
