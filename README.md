# cushion

[![Go Report Card](https://goreportcard.com/badge/github.com/ansurfen/cushion)](https://goreportcard.com/report/github.com/ansurfen/cushion)
![Software License](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat-square)
[![GoDoc](https://godoc.org/github.com/ansurfen/cushion?status.svg)](https://godoc.org/github.com/ansurfen/cushion)

English | [简体中文](./docs/zh_cn/README.md)

Cushion is a basic library for OpenCmd, Ark, Yock and Hulo...

## Features

* go-prompt: better powerful auto-completion with rgb color and key highlight and mode switch and so on
* utils: uniform interface about dynamic library (dll, dylib, so) and metatable (plist, regedit) and environment variable and so on
* runtime: dynamic language capabilities based on gopher-lua
* components: out-of-the-box components based on bubbletea

## Get start

To start, we'll fetch library using command.
```cmd
go get "github.com/ansurfen/cushion"
```

Then, import cushion from local repository into your project.
```go
package main

import (
    "github.com/ansurfen/cushion/utils"
    "github.com/ansurfen/cushion/runtime"
)
```