//go:build windows
// +build windows

package utils

import "syscall"

var (
	_ Plugin     = &WindowsPlugin{}
	_ PluginFunc = &WindowsPluginFunc{}
)

type WindowsPlugin struct {
	*syscall.LazyDLL
}

func NewPlugin(path string) (Plugin, error) {
	return &WindowsPlugin{
		LazyDLL: syscall.NewLazyDLL(path),
	}, nil
}

// Func return PluginFunc which is an abstract function to be exported dynamic library
// according to funcName. You can use PluginFunc to call function from dynamic library.
func (wp *WindowsPlugin) Func(plugin string) (PluginFunc, error) {
	return &WindowsPluginFunc{
		LazyProc: wp.LazyDLL.NewProc(plugin),
	}, nil
}

type WindowsPluginFunc struct {
	*syscall.LazyProc
}

// Call return excuted result from dynamic library
func (wpf *WindowsPluginFunc) Call(params ...uintptr) (uintptr, error) {
	ret, _, err := wpf.LazyProc.Call(params...)
	return ret, err
}
