package utils

import (
	"fmt"
	"testing"
)

type HuloEnv struct {
	Lang    string `yaml:"lang"`
	Patiner struct {
	} `yaml:"patiner"`
}

func (env *HuloEnv) SetLang(lang string) *HuloEnv {
	env.Lang = lang
	GetEnv().Commit("lang", env.Lang)
	return env
}

func (env *HuloEnv) Write() {
	GetEnv().Write()
}

func TestEnv(t *testing.T) {
	huloEnv := NewEnv(EnvOpt[HuloEnv]{
		Workdir: ".hulo",
		Subdirs: []string{"pkg", "loader"},
		BlankConf: `workdir: ""
lang: ""
patiner: ""
`,
		Payload: HuloEnv{},
	})
	fmt.Println(huloEnv)
	huloEnv.SetLang("zh_cn").Write()
}
