// +build !windows

package utils

import (
	"errors"
	"plugin"
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

func (ppf *PosixPluginFunc) Call(params ...uintptr) (uintptr, error) {
	switch f := ppf.Symbol.(type) {
	case func(...uintptr) uintptr:
		return f(params...), nil
	default:
		return uintptr(0), errors.New("fail to assert func")
	}
}
