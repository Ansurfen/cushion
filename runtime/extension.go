// +build extension

package runtime

import (
	"github.com/vadv/gopher-lua-libs/argparse"
	"github.com/vadv/gopher-lua-libs/aws/cloudwatch"
	"github.com/vadv/gopher-lua-libs/base64"
	"github.com/vadv/gopher-lua-libs/cert_util"
	"github.com/vadv/gopher-lua-libs/chef"
	"github.com/vadv/gopher-lua-libs/cmd"
	"github.com/vadv/gopher-lua-libs/crypto"
	"github.com/vadv/gopher-lua-libs/db"
	"github.com/vadv/gopher-lua-libs/filepath"
	"github.com/vadv/gopher-lua-libs/goos"
	"github.com/vadv/gopher-lua-libs/http"
	"github.com/vadv/gopher-lua-libs/humanize"
	"github.com/vadv/gopher-lua-libs/inspect"
	"github.com/vadv/gopher-lua-libs/ioutil"
	"github.com/vadv/gopher-lua-libs/json"
	"github.com/vadv/gopher-lua-libs/log"
	"github.com/vadv/gopher-lua-libs/pb"
	"github.com/vadv/gopher-lua-libs/plugin"
	"github.com/vadv/gopher-lua-libs/pprof"
	prometheus "github.com/vadv/gopher-lua-libs/prometheus/client"
	"github.com/vadv/gopher-lua-libs/regexp"
	"github.com/vadv/gopher-lua-libs/runtime"
	"github.com/vadv/gopher-lua-libs/shellescape"
	"github.com/vadv/gopher-lua-libs/stats"
	"github.com/vadv/gopher-lua-libs/storage"
	"github.com/vadv/gopher-lua-libs/strings"
	"github.com/vadv/gopher-lua-libs/tac"
	"github.com/vadv/gopher-lua-libs/tcp"
	"github.com/vadv/gopher-lua-libs/telegram"
	"github.com/vadv/gopher-lua-libs/template"
	"github.com/vadv/gopher-lua-libs/time"
	"github.com/vadv/gopher-lua-libs/xmlpath"
	"github.com/vadv/gopher-lua-libs/yaml"
)

func (vm *LuaVM) mountLibs() {
	vm.mat.Mount(LuaFuncs{
		"libs-plugin":      plugin.Loader,
		"libs-argparse":    argparse.Loader,
		"libs-base64":      base64.Loader,
		"libs-cert_util":   cert_util.Loader,
		"libs-chef":        chef.Loader,
		"libs-cloudwatch":  cloudwatch.Loader,
		"libs-cmd":         cmd.Loader,
		"libs-crypto":      crypto.Loader,
		"libs-db":          db.Loader,
		"libs-filepath":    filepath.Loader,
		"libs-goos":        goos.Loader,
		"libs-http":        http.Loader,
		"libs-humanize":    humanize.Loader,
		"libs-inspect":     inspect.Loader,
		"libs-ioutil":      ioutil.Loader,
		"libs-json":        json.Loader,
		"libs-log":         log.Loader,
		"libs-pb":          pb.Loader,
		"libs-pprof":       pprof.Loader,
		"libs-prometheus":  prometheus.Loader,
		"libs-regexp":      regexp.Loader,
		"libs-runtime":     runtime.Loader,
		"libs-shellescape": shellescape.Loader,
		"libs-stats":       stats.Loader,
		"libs-storage":     storage.Loader,
		"libs-strings":     strings.Loader,
		"libs-tac":         tac.Loader,
		"libs-tcp":         tcp.Loader,
		"libs-telegram":    telegram.Loader,
		"libs-template":    template.Loader,
		"libs-time":        time.Loader,
		"libs-xmlpath":     xmlpath.Loader,
		"libs-yaml":        yaml.Loader,
	})
}
