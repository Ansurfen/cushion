package utils

type Plugin interface {
	Func(string) (PluginFunc, error)
}

type PluginFunc interface {
	Call(...uintptr) (uintptr, error)
}
