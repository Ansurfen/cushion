//go:build windows
// +build windows

package utils

import (
	"errors"
	"os"
	"strconv"
	"strings"

	"golang.org/x/sys/windows/registry"
)

var _ EnvVar = &WinEnvVar{}

type WinEnvVar struct {
	self *WinEnv
	mode bool
}

func NewEnvVar() *WinEnvVar {
	return &WinEnvVar{
		self: NewWinEnv().ReadSysEnv().ReadUserVar(),
		mode: envVarUser,
	}
}

// windows: sys or user
func (env *WinEnvVar) SetPath(path string) error {
	switch path {
	case "sys":
		env.mode = envVarSys
	case "user":
		env.mode = envVarUser
	default:
		return errors.New("err mode")
	}
	return nil
}

func (env *WinEnvVar) Set(k string, v any) error {
	var regVal RegistryValue
	switch vv := v.(type) {
	case string:
		regVal = NewRegistryValue(registry.SZ, vv)
	case []string:
		regVal = NewRegistryValue(registry.EXPAND_SZ, strings.Join(vv, ";"))
	case int:
		regVal = NewRegistryValue(registry.DWORD, strconv.Itoa(vv))
	default:
		return errors.New("invalid value")
	}
	if !env.mode {
		env.self.SetUserEnv(k, regVal)
	} else {
		env.self.SetSysEnv(k, regVal)
	}
	return nil
}

func (env *WinEnvVar) SafeSet(k string, v any) error {
	var regVal RegistryValue
	switch vv := v.(type) {
	case string:
		regVal = NewRegistryValue(registry.SZ, vv)
	case []string:
		regVal = NewRegistryValue(registry.EXPAND_SZ, strings.Join(vv, ";"))
	case int:
		regVal = NewRegistryValue(registry.DWORD, strconv.Itoa(vv))
	default:
		return errors.New("invalid value")
	}
	if !env.mode {
		env.self.SafeSetUserEnv(k, regVal)
	} else {
		env.self.SafeSetSysEnv(k, regVal)
	}
	return nil
}

func (env *WinEnvVar) Unset(k string) error {
	opt := EnvVarDeleteOpt{
		Rules: []string{k},
		Safe:  false,
	}
	if !env.mode {
		env.self.DeleteUserVar(opt)
	} else {
		env.self.DeleteSysVar(opt)
	}
	return nil
}

func (env *WinEnvVar) SetL(k, v string) error {
	return os.Setenv(k, v)
}

func (env *WinEnvVar) SafeSetL(k, v string) error {
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

func (env *WinEnvVar) Export(file string) error {
	opt := EnvVarExportOpt{
		File: file,
	}
	if !env.mode {
		env.self.ExportUserVar(opt)
	} else {
		env.self.ExportSysVar(opt)
	}
	return nil
}

func (env *WinEnvVar) Load(opt EnvVarLoadOpt) error {
	env.self.LoadEnvVar(WinEnvVarLoadOpt{
		File: opt.file,
		Spec: env.mode,
	})
	return nil
}

func (env *WinEnvVar) Print() {
	if !env.mode {
		env.self.DumpUserVar()
	} else {
		env.self.DumpSysVar()
	}
}
