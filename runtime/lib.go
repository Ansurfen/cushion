package runtime

import (
	"path"
	"regexp"
	"runtime"
	"strings"

	"github.com/ansurfen/cushion/components"
	"github.com/ansurfen/cushion/utils"
	"github.com/ansurfen/cushion/utils/build"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
	lua "github.com/yuin/gopher-lua"
	luar "layeh.com/gopher-luar"
)

func loadVM(lvm *lua.LState) int {
	return LuaModuleLoader(lvm, LuaFuncs{
		"Eval":      vmEval,
		"EvalFile":  vmEvalFile,
		"SetGlobal": vmSetGlobal,
		"Export":    vmExport,
		"OS":        vmOS,
		"Arch":      vmArch,
		"Workdir":   vmWorkdir,
	})
}

func vmWorkdir(lvm *lua.LState) int {
	lvm.Push(lua.LString(utils.GetEnv().Workdir()))
	return 1
}

func vmOS(lvm *lua.LState) int {
	lvm.Push(lua.LString(runtime.GOOS))
	return 1
}

func vmArch(lvm *lua.LState) int {
	lvm.Push(lua.LString(runtime.GOARCH))
	return 1
}

func vmExport(lvm *lua.LState) int {
	key := utils.FirstUpper(utils.RandString(8))
	value := lvm.CheckString(1)
	tbl := lvm.CheckAny(2)
	lvm.SetGlobal(key, tbl)
	GetExport().Set(key, value)
	return 0
}

func vmEval(lvm *lua.LState) int {
	err := lvm.DoString(lvm.CheckString(1))
	errHandle(lvm, err)
	return 1
}

func vmEvalFile(lvm *lua.LState) int {
	err := lvm.DoFile(lvm.CheckString(1))
	errHandle(lvm, err)
	return 1
}

func vmSetGlobal(lvm *lua.LState) int {
	lvm.SetGlobal(lvm.CheckString(1), lvm.CheckAny(2))
	return 0
}

func loadTui(lvm *lua.LState) int {
	return LuaModuleLoader(lvm, LuaFuncs{
		"FancyList":    tuiFancyList,
		"SimpleList":   tuiSimpleList,
		"Spinner":      tuiSpinner,
		"TextInput":    tuiTextInput,
		"MultiSelect":  tuiMultiSelect,
		"BatchSpinner": tuiBatchSpinner,
	})
}

func tuiSpinner(lvm *lua.LState) int {
	components.UseSpinner(components.DefaultSpinnerStyle(), &components.SpinnerPayLoad{
		Callback: func() {
			LuaDoFunc(lvm, lvm.CheckFunction(1))
		},
	})
	return 0
}

func tuiFancyList(lvm *lua.LState) int {
	items := []list.Item{}
	lvm.CheckTable(1).ForEach(func(idx, item lua.LValue) {
		items = append(items, components.FancyListItem{
			ChoiceTitle:  item.(*lua.LTable).RawGetString("title").String(),
			ChoiceDetail: item.(*lua.LTable).RawGetString("detail").String(),
		})
	})
	choice := components.UseFancyList(components.DefaultFancyListStyle(), &components.FancyListPayLoad{
		Title:   "",
		Choices: items,
	})
	lvm.Push(luar.New(lvm, choice))
	return 1
}

func tuiSimpleList(lvm *lua.LState) int {
	payload := &components.SimpleListPayLoad{}
	lvm.CheckTable(1).ForEach(func(l1, l2 lua.LValue) {
		switch l1.String() {
		case "title":
			payload.Title = l2.String()
		case "choices":
			l2.(*lua.LTable).ForEach(func(l1, l2 lua.LValue) {
				payload.Choices = append(payload.Choices, l2.String())
			})
		}
	})
	style := components.DefaultSimpleListStyle()
	style.TitleStyle = lipgloss.NewStyle()
	out := components.UseSimpleList(style, payload)
	lvm.Push(lua.LNumber(out))
	return 1
}

func tuiTextInput(lvm *lua.LState) int {
	texts := []components.TextInputFormat{}
	lvm.CheckTable(1).ForEach(func(idx, tbl lua.LValue) {
		text := tbl.(*lua.LTable)
		echomode := false
		switch text.RawGetString("echomode").String() {
		case "true":
			echomode = true
		default:
		}
		texts = append(texts, components.TextInputFormat{
			Name:     text.RawGetString("name").String(),
			EchoMode: echomode,
		})
	})
	if len(texts) > 0 {
		res := components.UseTextInput(components.DefaultTextStyle(), &components.TextInputPayLoad{
			Texts: texts,
		})
		tbl := &lua.LTable{}
		for _, r := range res {
			tbl.Append(lua.LString(r))
		}
		lvm.Push(tbl)
		return 1
	}
	return 0
}

func tuiMultiSelect(lvm *lua.LState) int {
	tbl := lvm.CheckTable(1)
	title := tbl.RawGetString("title").String()
	choices := []string{}
	tbl.RawGetString("choices").(*lua.LTable).ForEach(func(idx, choice lua.LValue) {
		choices = append(choices, choice.String())
	})
	if len(choices) > 0 {
		res := components.UseMultiSelect(components.DefaultMultiSelectStyle(), &components.MultiSelectPayLoad{
			Title:   title,
			Choices: choices,
		})
		tbl := &lua.LTable{}
		for _, r := range res {
			tbl.Append(lua.LNumber(r))
		}
		lvm.Push(tbl)
		return 1
	}
	return 0
}

func tuiBatchSpinner(lvm *lua.LState) int {
	tasks := []components.BatchTask{}
	lvm.CheckTable(1).ForEach(func(idx, tbl lua.LValue) {
		tasks = append(tasks, components.BatchTask{
			Name: tbl.(*lua.LTable).RawGetString("name").String(),
			Callback: func() bool {
				if err := LuaDoFunc(lvm, tbl.(*lua.LTable).RawGetString("callback").(*lua.LFunction)); err != nil {
					return false
				}
				return true
			},
		})
	})
	components.UseBatchSpinner(components.DefaultBatchSpinnerStyle(), &components.BatchSpinnerPayLoad{
		Task: tasks,
	})
	return 1
}

func LoadCheck(lvm *lua.LState) int {
	return LuaModuleLoader(lvm, LuaFuncs{
		"CheckEnv":      checkEnv,
		"CheckVersion":  checkVersion,
		"FormatVersion": checkFormatVersion,
	})
}

func checkVersion(lvm *lua.LState) int {
	want := utils.NewCheckedVersion(lvm.CheckString(1))
	got := utils.NewCheckedVersion(lvm.CheckString(2))
	if want.Compare(got) {
		lvm.Push(lua.LTrue)
	} else {
		lvm.Push(lua.LFalse)
	}
	return 1
}

func checkFormatVersion(lvm *lua.LState) int {
	rawVersion := lvm.CheckString(1)
	targetCnt := lvm.CheckInt(2)
	cnt := 0
	curVersion := ""
	for _, ch := range rawVersion {
		if targetCnt == cnt-1 {
			break
		}
		if ch == '.' {
			cnt++
		}
		curVersion += string(ch)
	}
	if curVersionLen := len(curVersion); curVersion[curVersionLen-1] == '.' {
		curVersion = curVersion[:curVersionLen-1]
	}
	lvm.Push(lua.LString(curVersion))
	return 1
}

func checkEnv(lvm *lua.LState) int {
	cmd := lvm.CheckString(1)
	out, err := utils.Exec(cmd)
	res := ""
	if err != nil {
		res = ""
	} else {
		filter := lvm.CheckString(2)
		reg := regexp.MustCompile(filter)
		regRes := reg.FindStringSubmatch(string(out))
		if len(regRes) > 1 {
			res = regRes[1]
		}
	}
	lvm.Push(lua.LString(res))
	return 1
}

func loadIO(lvm *lua.LState) int {
	return LuaModuleLoader(lvm, LuaFuncs{
		"Fetch":  ioFetch,
		"Unzip":  ioUnzip,
		"Exec":   ioExec,
		"Mkdirs": ioMkdirs,
		"Printf": ioPrintf,
	})
}

func ioPrintf(lvm *lua.LState) int {
	title := []string{}
	rows := [][]string{}
	lvm.CheckTable(1).ForEach(func(idx, el lua.LValue) {
		title = append(title, el.String())
	})
	lvm.CheckTable(2).ForEach(func(ri, row lua.LValue) {
		tmp := []string{}
		row.(*lua.LTable).ForEach(func(fi, field lua.LValue) {
			tmp = append(tmp, field.String())
		})
		rows = append(rows, tmp)
	})
	utils.Prinf(utils.PrintfOpt{MaxLen: 10}, title, rows)
	return 0
}

func ioMkdirs(lvm *lua.LState) int {
	err := utils.SafeMkdirs(lvm.CheckString(1))
	errHandle(lvm, err)
	return 1
}

func ioExec(lvm *lua.LState) int {
	out, err := utils.ExecStr(lvm.CheckString(1))
	lvm.Push(lua.LString(utils.ConvertByte2String(out, utils.GB18030)))
	errHandle(lvm, err)
	return 2
}

func ioFetch(lvm *lua.LState) int {
	_, err := utils.FetchFile(lvm.CheckString(1), lvm.CheckString(2))
	if err != nil {
		lvm.Push(lua.LString(err.Error()))
	} else {
		lvm.Push(lua.LNil)
	}
	return 1
}

func ioUnzip(lvm *lua.LState) int {
	err := utils.Unzip(lvm.CheckString(1), lvm.CheckString(2))
	errHandle(lvm, err)
	return 1
}

func loadTmpl(lvm *lua.LState) int {
	return LuaModuleLoader(lvm, LuaFuncs{
		"Compile": tmplCompile,
	})
}

func tmplCompile(lvm *lua.LState) int {
	dict := make(map[string]string)
	lvm.CheckTable(2).ForEach(func(l1, l2 lua.LValue) {
		dict[l1.String()] = l2.String()
	})
	tmpl := build.NewTemplate()
	out, err := tmpl.OnceParse(lvm.CheckString(1), dict)
	lvm.Push(lua.LString(out))
	lvm.Push(luar.New(lvm, err))
	return 2
}

func loadCrypto(lvm *lua.LState) int {
	return LuaModuleLoader(lvm, LuaFuncs{
		"MD5":    cryptoMD5,
		"SHA256": cryptoSHA256,
	})
}

func cryptoMD5(lvm *lua.LState) int {
	lvm.Push(lua.LString(utils.MD5(lvm.CheckString(1))))
	return 1
}

func cryptoSHA256(lvm *lua.LState) int {
	lvm.Push(lua.LString(utils.SHA256(lvm.CheckString(1))))
	return 1
}

func loadTime(lvm *lua.LState) int {
	return LuaModuleLoader(lvm, LuaFuncs{
		"Now": timeNow,
	})
}

func timeNow(lvm *lua.LState) int {
	lvm.Push(lua.LNumber(utils.NowTimestamp()))
	return 1
}

func loadPath(lvm *lua.LState) int {
	return LuaModuleLoader(lvm, LuaFuncs{
		"IsAbs":    pathIsAbs,
		"Base":     pathBase,
		"Ext":      pathExt,
		"Clean":    pathClean,
		"Dir":      pathDir,
		"Join":     pathJoin,
		"Split":    pathSplit,
		"Match":    pathMatch,
		"Filename": pathFilename,
	})
}

func pathIsAbs(lvm *lua.LState) int {
	lvm.Push(lua.LBool(path.IsAbs(lvm.CheckString(1))))
	return 1
}

func pathBase(lvm *lua.LState) int {
	lvm.Push(lua.LString(path.Base(lvm.CheckString(1))))
	return 1
}

func pathExt(lvm *lua.LState) int {
	lvm.Push(lua.LString(path.Ext(lvm.CheckString(1))))
	return 1
}

func pathClean(lvm *lua.LState) int {
	lvm.Push(lua.LString(path.Clean(lvm.CheckString(1))))
	return 1
}

func pathDir(lvm *lua.LState) int {
	lvm.Push(lua.LString(path.Dir(lvm.CheckString(1))))
	return 1
}

func pathJoin(lvm *lua.LState) int {
	elem := []string{}
	for i := 1; i <= lvm.GetTop(); i++ {
		elem = append(elem, lvm.CheckString(i))
	}
	lvm.Push(lua.LString(path.Join(elem...)))
	return 1
}

func pathSplit(lvm *lua.LState) int {
	dir, file := path.Split(lvm.CheckString(1))
	lvm.Push(lua.LString(dir))
	lvm.Push(lua.LString(file))
	return 2
}

func pathMatch(lvm *lua.LState) int {
	pattern := lvm.CheckString(1)
	name := lvm.CheckString(2)
	matched, err := path.Match(pattern, name)
	lvm.Push(lua.LBool(matched))
	errHandle(lvm, err)
	return 2
}

func pathFilename(lvm *lua.LState) int {
	fullpath := lvm.CheckString(1)
	lvm.Push(lua.LString(utils.Filename(fullpath)))
	return 1
}

func loadStrings(lvm *lua.LState) int {
	return LuaModuleLoader(lvm, LuaFuncs{
		"Split":     stringsSplit,
		"Cut":       stringsCut,
		"Contains":  stringsContains,
		"HasPrefix": stringsHasPrefix,
		"HasSuffix": stringsHasSuffix,
	})
}

func stringsSplit(lvm *lua.LState) int {
	slice := strings.Split(lvm.CheckString(1), lvm.CheckString(2))
	lvm.Push(strSlice2Table(slice))
	return 1
}

func stringsCut(lvm *lua.LState) int {
	before, after, ok := strings.Cut(lvm.CheckString(1), lvm.CheckString(2))
	lvm.Push(lua.LString(before))
	lvm.Push(lua.LString(after))
	lvm.Push(lua.LBool(ok))
	return 3
}

func stringsContains(lvm *lua.LState) int {
	boolHandle(lvm, strings.Contains(lvm.CheckString(1), lvm.CheckString(2)))
	return 1
}

func stringsHasPrefix(lvm *lua.LState) int {
	boolHandle(lvm, strings.HasPrefix(lvm.CheckString(1), lvm.CheckString(2)))
	return 1
}

func stringsHasSuffix(lvm *lua.LState) int {
	boolHandle(lvm, strings.HasSuffix(lvm.CheckString(1), lvm.CheckString(2)))
	return 1
}
