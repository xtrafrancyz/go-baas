# go-baas

Bcrypt as a service on Golang.

Inspired by [node-baas](https://github.com/auth0/node-baas). This implementation does **not** include
client for service, metrics, own protocol or any other enterprise features.

## Api

There is only 2 http methods:

```
GET /hash?raw=PASSWORD
$2a$10$8WN5TKy2oTVbKrPhvA.lX.n2ef9eKmGd0.QXekflZBwd6LKdBiJ.C
```
```
GET /verify?raw=PASSWORD&hash=HASH
OK or FAIL
```

## Setup
```
go get github.com/xtrafrancyz/go-baas
go build
./go-baas -bind 127.0.0.1:8085 -threads 1 -cost 10
```
