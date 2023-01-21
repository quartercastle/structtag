# structtag

[![Version](https://img.shields.io/github/release/quartercastle/structtag.svg)](https://github.com/quartercastle/structtag/releases)
[![GoDoc](https://godoc.org/github.com/quartercastle/tag?status.svg)](https://godoc.org/github.com/quartercastle/structtag)
[![Go Report Card](https://goreportcard.com/badge/github.com/quartercastle/structtag)](https://goreportcard.com/report/github.com/quartercastle/structtag)


The motivation behind this package is that the [`StructTag`](https://github.com/golang/go/blob/0377f061687771eddfe8de78d6c40e17d6b21a39/src/reflect/type.go#L1110)
implementation shipped with Go's standard library is very limited in
detecting a malformed StructTag and each time `StructTag.Get(key)` gets called,
it results in the `StructTag` being parsed again. 
This package provides a way to parse the `StructTag` into a `structtag.Map`.

```go
// Example of struct using tags to append metadata to fields.
type Server struct {
	Host string `json:"host" env:"SERVER_HOST" default:"localhost"`
	Port int    `json:"port" env:"SERVER_PORT" default:"3000"`
}
```

### Install
```
go get github.com/quartercastle/structtag
```

### Usage
```go
tags, err := structtag.Parse(`json:"host" env:"SERVER_HOST"`)

if err != nil {
  panic(err)
}

fmt.Println(tags["json"])
```
See [godoc](https://godoc.org/github.com/quartercastle/structtag) for full documentation.

### License
This project is licensed under the [MIT License](https://github.com/quartercastle/structtag/blob/master/LICENSE).
