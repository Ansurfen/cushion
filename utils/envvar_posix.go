//go:build !windows
// +build !windows

package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
)

var _ EnvVar = &PosixEnvVar{}

type PosixEnvVar struct{}

func NewEnvVar() *PosixEnvVar {
	return &PosixEnvVar{}
}

// SetPath set operate target: posix: /etc/enviroment, this only is empty method.
func (env *PosixEnvVar) SetPath(path string) error { return nil }

// set global enviroment variable
func (env *PosixEnvVar) Set(k string, v any) error {
	if _, err := Exec("export", fmt.Sprintf("%s=%s", k, v)); err != nil {
		return err
	}
	return nil
}

// set global enviroment variable when key isn't exist
func (env *PosixEnvVar) SafeSet(k string, v any) error {
	vv, err := Exec("echo", "$"+k)
	if err != nil {
		return err
	}
	if string(vv) != "\n" {
		if _, err := Exec("export", fmt.Sprintf("%s=%s", k, v)); err != nil {
			return err
		}
	}
	return nil
}

// set local enviroment variable
func (env *PosixEnvVar) SetL(k, v string) error {
	return os.Setenv(k, v)
}

// set local enviroment variable when key isn't exist
func (env *PosixEnvVar) SafeSetL(k, v string) error {
	exist := false
	for _, e := range os.Environ() {
		if kk, _, ok := strings.Cut(e, "="); ok && kk == k {
			exist = true
			break
		}
	}
	if !exist {
		return os.Setenv(k, v)
	}
	return errors.New("var exist already")
}

// unset (delete) global enviroment variable
func (env *PosixEnvVar) Unset(k string) error {
	if _, err := Exec("unset", k); err != nil {
		return err
	}
	return nil
}

// export current enviroment string into specify file
func (env *PosixEnvVar) Export(file string) error {
	dict := make(map[string]string)
	for _, e := range os.Environ() {
		if k, v, ok := strings.Cut(e, "="); ok {
			dict[k] = v
		}
	}
	raw, err := json.Marshal(dict)
	if err != nil {
		return err
	}
	return WriteFile(file, raw)
}

// load exported env from disk
func (env *PosixEnvVar) Load(opt EnvVarLoadOpt) error {
	raw, err := ReadStraemFromFile(opt.file)
	if err != nil {
		return err
	}
	dict := make(map[string]string)
	if json.Unmarshal(raw, &dict) != nil {
		return err
	}
	for _, k := range opt.keys {
		if v, ok := dict[k]; ok {
			if opt.safe {
				env.SafeSet(k, v)
			} else {
				env.Set(k, v)
			}
		}
	}
	return nil
}

// Print enviroment variable
func (env *PosixEnvVar) Print() {
	for _, e := range os.Environ() {
		if k, v, ok := strings.Cut(e, "="); ok {
			fmt.Printf("%s: %s", k, v)
		}
	}
}
