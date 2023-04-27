package utils

// EnvVar is an interface to abstract different os enviroment variable
type EnvVar interface {
	// SetPath set operate target:
	// windows: sys or user, it's required in windows.
	// posix: /etc/enviroment, this only is empty method.
	SetPath(string) error

	// set global enviroment variable
	Set(string, any) error
	// set global enviroment variable when key isn't exist
	SafeSet(string, any) error

	// unset (delete) global enviroment variable
	Unset(string) error

	// set local enviroment variable
	SetL(string, string) error
	// set local enviroment variable when key isn't exist
	SafeSetL(string, string) error

	// export current enviroment string into specify file
	Export(string) error
	// load enviroment string to be export from disk
	Load(EnvVarLoadOpt) error

	// Print current enviroment variable
	Print()
}

type EnvVarLoadOpt struct {
	file string
	keys []string
	safe bool
}
