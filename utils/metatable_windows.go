//go:build windows
// +build windows

package utils

import (
	"path/filepath"
	"strings"

	"golang.org/x/sys/windows/registry"
)

type WinMetaTable struct {
	page *RegistryPage
	sub  []*RegistryPage
}

func CreateMetaTable(path string) (MetaTable, error) {
	tbl := &WinMetaTable{
		sub: make([]*RegistryPage, 0),
	}
	path = filepath.Clean(path)
	if root, path, ok := strings.Cut(path, "\\"); ok {
		switch root {
		case "ROOT":
			tbl.page = CreateRegistryPage(registry.CLASSES_ROOT, path)
		case "USER":
			tbl.page = CreateRegistryPage(registry.CURRENT_USER, path)
		case "LOCAL_MACHINE":
			tbl.page = CreateRegistryPage(registry.LOCAL_MACHINE, path)
		case "USERS":
			tbl.page = CreateRegistryPage(registry.USERS, path)
		case "CURRENT_CONFIG":
			tbl.page = CreateRegistryPage(registry.CURRENT_CONFIG, path)
		}
	}
	return tbl, nil
}

func OpenMetaTable(path string) (MetaTable, error) {
	tbl := &WinMetaTable{
		sub: make([]*RegistryPage, 0),
	}
	path = filepath.Clean(path)
	if root, path, ok := strings.Cut(path, "\\"); ok {
		switch root {
		case "ROOT":
			tbl.page = OpenRegistryPage(registry.CLASSES_ROOT, path)
		case "USER":
			tbl.page = OpenRegistryPage(registry.CURRENT_USER, path)
		case "LOCAL_MACHINE":
			tbl.page = OpenRegistryPage(registry.LOCAL_MACHINE, path)
		case "USERS":
			tbl.page = OpenRegistryPage(registry.USERS, path)
		case "CURRENT_CONFIG":
			tbl.page = OpenRegistryPage(registry.CURRENT_CONFIG, path)
		}
	}
	return tbl, nil
}

func (tbl *WinMetaTable) GetValue(v string) MetaValue {
	return GetValue(tbl.page.key, v)
}

func (tbl *WinMetaTable) SetValue(v MetaValue) {
	switch vv := v.(type) {
	case MetaMap:
		for name, value := range vv {
			switch vv := value.(type) {
			case int:
				tbl.page.SetValue(name, DWordValue{
					val: uint64(vv),
				})
			case string:
				tbl.page.SetValue(name, SZValue{
					val: vv,
				})
			case []string:
				tbl.page.SetValue(name, ExpandSZValue{
					val: vv,
				})
			default:

			}
		}
	case MetaArr:
	default:
	}
}

func (tbl *WinMetaTable) SafeSetValue(v MetaValue) {
	switch vv := v.(type) {
	case MetaMap:
		for name, value := range vv {
			switch vv := value.(type) {
			case int:
				tbl.page.SafeSetValue(name, DWordValue{
					val: uint64(vv),
				})
			case string:
				tbl.page.SafeSetValue(name, SZValue{
					val: vv,
				})
			case []string:
				tbl.page.SafeSetValue(name, ExpandSZValue{
					val: vv,
				})
			default:

			}
		}
	case MetaArr:
	default:
	}
}

func (tbl *WinMetaTable) CreateSubTable(name string) MetaTable {
	child := tbl.page.CreateSubKey(name)
	tbl.sub = append(tbl.sub, child)
	return &WinMetaTable{
		page: child,
	}
}

func (tbl *WinMetaTable) Write() error {
	return nil
}

func (tbl *WinMetaTable) Backup() error {
	return tbl.page.Backup()
}

func (tbl *WinMetaTable) Close() {
	for _, child := range tbl.sub {
		child.Free()
	}
	tbl.page.Free()
}
