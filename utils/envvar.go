package utils

type EnvVar interface {
	// windows: sys or user
	// linux: /etc/enviroment
	SetPath(string) error

	// global env variable
	Set(string, any) error
	SafeSet(string, any) error

	Unset(string) error

	// local env variable
	SetL(string, string) error
	SafeSetL(string, string) error

	// export current env string
	Export(string) error
	// load exported env from disk
	Load(EnvVarLoadOpt) error

	Print()
}

type EnvVarLoadOpt struct {
	file string
	keys []string
	safe bool
}
