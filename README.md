# go-baas

Bcrypt as a service on Golang.

Inspired by [node-baas](https://github.com/auth0/node-baas). This implementation does **not** include
client for service, metrics, own protocol or any other enterprise features.

## API

There is only 2 http methods:

```
POST /hash
Content-Type: application/x-www-form-urlencoded

raw=PASSWORD&cost=10

// $2a$10$8WN5TKy2oTVbKrPhvA.lX.n2ef9eKmGd0.QXekflZBwd6LKdBiJ.C
```

```
POST /verify
Content-Type: application/x-www-form-urlencoded

raw=PASSWORD&hash=HASH

// OK or FAIL
```

## Installation

```
go install github.com/xtrafrancyz/go-baas@latest
```

## Running

```
go-baas -bind 127.0.0.1:8085 -threads 1 -cost 10
```
