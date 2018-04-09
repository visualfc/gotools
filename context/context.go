package context

import (
	"go/build"
	"os"
)

var (
	fnLookupEnv = os.LookupEnv
)

func SetLookupEnv(fn func(key string) (string, bool)) {
	if fn != nil {
		fnLookupEnv = fn
	} else {
		fnLookupEnv = os.LookupEnv
	}
}

func Default() *build.Context {
	return &build.Default
}

func System() *build.Context {
	return NewContext(fnLookupEnv)
}

func NewContext(env func(key string) (string, bool)) *build.Context {
	c := build.Default
	if v, ok := env("GOARCH"); ok {
		c.GOARCH = v
	}
	if v, ok := env("GOOS"); ok {
		c.GOOS = v
	}
	if v, ok := env("GOROOT"); ok {
		c.GOROOT = v
	}
	if v, ok := env("GOPATH"); ok {
		c.GOPATH = v
	}
	if v, ok := env("CGO_ENABLED"); ok {
		switch v {
		case "1":
			c.CgoEnabled = true
		case "0":
			c.CgoEnabled = false
		}
	}
	return &c
}
