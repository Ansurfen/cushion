//go:build windows
// +build windows

package utils

import "syscall"

type WindowsPlugin struct {
	*syscall.LazyDLL
}

func NewPlugin(path string) (Plugin, error) {
	return &WindowsPlugin{
		LazyDLL: syscall.NewLazyDLL(path),
	}, nil
}

func (wp *WindowsPlugin) Func(plugin string) (PluginFunc, error) {
	return &WindowsPluginFunc{
		LazyProc: wp.LazyDLL.NewProc(plugin),
	}, nil
}

type WindowsPluginFunc struct {
	*syscall.LazyProc
}

func (wpf *WindowsPluginFunc) Call(params ...uintptr) (uintptr, error) {
	ret, _, err := wpf.LazyProc.Call(params...)
	return ret, err
}
