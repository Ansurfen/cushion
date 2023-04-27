package utils

// Plugin is an interface to abstract dynamic library (dll, dylib, so)
type Plugin interface {
	// Func return PluginFunc which is an abstract function to be exported dynamic library
	// according to funcName. You can use PluginFunc to call function from dynamic library.
	Func(string) (PluginFunc, error)
}

// PluginFunc is an interface to abstract function to be exported dynamic library
type PluginFunc interface {
	// Call return excuted result from dynamic library
	Call(...uintptr) (uintptr, error)
}
