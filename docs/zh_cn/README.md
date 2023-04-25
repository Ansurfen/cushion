# cushion

[![Go Report Card](https://goreportcard.com/badge/github.com/ansurfen/cushion)](https://goreportcard.com/report/github.com/ansurfen/cushion)
![Software License](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat-square)
[![GoDoc](https://godoc.org/github.com/ansurfen/cushion?status.svg)](https://godoc.org/github.com/ansurfen/cushion)

[English](../../README.md) | 简体中文

Cushion 是 OpenCmd、Ark、Hulo、Yock 的基础库

## 特性

* go-prompt: 更强大的自动补全，支持RGB颜色、关键字高亮以及模式选择...
* utils: 统一的动态库（dll, dylib, so）封装、系统及软件信息元数据表封装（plist, regedit）、环境变量...
* runtime: 具备动态语言的能力基于gopher-lua封装
* components: 基于bubbletea封装，开箱即用的组件库

## 开始

最开始，我们将利用命令获取库
```cmd
go get "github.com/ansurfen/cushion"
```

然后，导入本地的库进入项目
```go
package main

import (
    "github.com/ansurfen/cushion/utils"
    "github.com/ansurfen/cushion/runtime" // 
)
```