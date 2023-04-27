//go:build !windows
// +build !windows

package utils

import (
	"errors"
	"plugin"
)

var (
	_ Plugin     = &PosixPlugin{}
	_ PluginFunc = &PosixPluginFunc{}
)

type PosixPlugin struct {
	*plugin.Plugin
}

func NewPlugin(path string) (Plugin, error) {
	plugin, err := plugin.Open(path)
	return &PosixPlugin{
		Plugin: plugin,
	}, err
}

// Func return PluginFunc which is an abstract function to be exported dynamic library
// according to funcName. You can use PluginFunc to call function from dynamic library.
func (pp *PosixPlugin) Func(name string) (PluginFunc, error) {
	sym, err := pp.Lookup(name)
	if err != nil {
		return nil, err
	}
	return &PosixPluginFunc{
		Symbol: sym,
	}, nil
}

type PosixPluginFunc struct {
	plugin.Symbol
}

// Call return excuted result from dynamic library
func (ppf *PosixPluginFunc) Call(params ...uintptr) (uintptr, error) {
	switch f := ppf.Symbol.(type) {
	case func(...uintptr) uintptr:
		return f(params...), nil
	default:
		return uintptr(0), errors.New("fail to assert func")
	}
}
