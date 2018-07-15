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

## Description
Trade-offs:
1. only Sydney is supported and only wind and temperature are returned 
2. it's still possible to burden service providers once cache expired
3. no authorization is performed
4. it uses in memory cache which is not scalable

Things could be improved:
All the points above could be considered if time permitted. For example, I'm interested to use memcached to serve as cache so it's more scalable. Also, I want to build a simple front end for the service using React.