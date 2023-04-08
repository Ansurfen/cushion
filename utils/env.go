package utils

import (
	"fmt"
	"os/user"
	"path"
	"path/filepath"

	"github.com/spf13/viper"
)

var env *Env

const (
	BlankConf = `workdir: ""
`
	ConfFile = "conf.yaml"
)

func GetEnv() *Env {
	return env
}

type Env struct {
	workdir string `yaml:"workdir"`
	user    *user.User
	conf    *viper.Viper
	file    string
}

func InitEnv(workdir string, subdirs []string) {
	env = NewEnv()
	env.workdir = filepath.ToSlash(path.Join(env.workdir, workdir))
	if ok, err := PathIsExist(env.workdir); err != nil {
		panic(err)
	} else if !ok {
		if err := env.initWorkspace(subdirs); err != nil {
			panic(err)
		}
	}
	env.file = path.Join(env.workdir, ConfFile)
	if ok, err := PathIsExist(env.file); err != nil {
		panic(err)
	} else if ok {
		env.Read(env.file)
	} else {
		err := WriteFile(env.file, []byte(BlankConf))
		if err != nil {
			panic(err)
		}
		env.Read(env.file)
	}
}

func NewEnv() *Env {
	curUser, err := user.Current()
	if err != nil {
		panic(err)
	}
	return &Env{
		user:    curUser,
		workdir: curUser.HomeDir,
	}
}

func (env *Env) Workdir() string {
	return env.workdir
}

func (env *Env) initWorkspace(dirs []string) error {
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

func (env *Env) Dump() {
	fmt.Println("workdir: ", env.workdir)
}

func (env *Env) Read(path string) {
	env.conf = NewConfFromPath(path)
	if wd := env.conf.GetString("workdir"); len(wd) > 0 {
		env.workdir = filepath.ToSlash(wd)
	}
	if err := env.conf.Unmarshal(env); err != nil {
		panic(err)
	}
}

func (env *Env) Write() {
	if err := env.conf.WriteConfig(); err != nil {
		panic(err)
	}
}
