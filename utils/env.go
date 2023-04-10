package utils

import (
	"fmt"
	"os/user"
	"path"
	"path/filepath"

	"github.com/spf13/viper"
)

var env *BaseEnv

const (
	BlankConf = `workdir: ""`
	ConfFile  = "conf.yaml"
)

func GetEnv() *BaseEnv {
	return env
}

type BaseEnv struct {
	workdir string `yaml:"workdir"`
	user    *user.User
	conf    *viper.Viper
	file    string
}

type EnvOpt[T any] struct {
	Payload   T
	Workdir   string
	Subdirs   []string
	BlankConf string
}

func NewEnv[T any](opt EnvOpt[T]) *T {
	env = NewBaseEnv()
	env.workdir = filepath.ToSlash(path.Join(env.workdir, opt.Workdir))
	if ok, err := PathIsExist(env.workdir); err != nil {
		panic(err)
	} else if !ok {
		if err := env.initWorkspace(opt.Subdirs); err != nil {
			panic(err)
		}
	}
	env.file = path.Join(env.workdir, ConfFile)
	if ok, err := PathIsExist(env.file); err != nil {
		panic(err)
	} else if ok {
		env.ReadWithBind(env.file, &opt.Payload)
	} else {
		bc := opt.BlankConf
		if len(bc) == 0 {
			bc = BlankConf
		}
		err := WriteFile(env.file, []byte(bc))
		if err != nil {
			panic(err)
		}
		env.ReadWithBind(env.file, &opt.Payload)
	}
	return &opt.Payload
}

func NewBaseEnv() *BaseEnv {
	curUser, err := user.Current()
	if err != nil {
		panic(err)
	}
	return &BaseEnv{
		user:    curUser,
		workdir: curUser.HomeDir,
	}
}

func (env *BaseEnv) Workdir() string {
	return env.workdir
}

func (env *BaseEnv) initWorkspace(dirs []string) error {
	workdir := env.workdir
	if err := SafeMkdirs(workdir); err != nil {
		return err
	}
	for _, dir := range dirs {
		dir = path.Join(workdir, dir)
		if err := SafeMkdirs(dir); err != nil {
			return err
		}
	}
	return nil
}

func (env *BaseEnv) Dump() {
	fmt.Println("workdir: ", env.workdir)
}

func (env *BaseEnv) Read(path string) {
	env.conf = NewConfFromPath(path)
	if wd := env.conf.GetString("workdir"); len(wd) > 0 {
		env.workdir = filepath.ToSlash(wd)
	}
	if err := env.conf.Unmarshal(env); err != nil {
		panic(err)
	}
}

func (env *BaseEnv) ReadWithBind(path string, payload any) {
	env.conf = NewConfFromPath(path)
	if wd := env.conf.GetString("workdir"); len(wd) > 0 {
		env.workdir = filepath.ToSlash(wd)
	}
	if err := env.conf.Unmarshal(payload); err != nil {
		panic(err)
	}
}

func (env *BaseEnv) Commit(key string, value any) *BaseEnv {
	env.conf.Set(key, value)
	return env
}

func (env *BaseEnv) Write() {
	if err := env.conf.WriteConfig(); err != nil {
		panic(err)
	}
}
