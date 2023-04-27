//go:build !windows
// +build !windows

package utils

import (
	"strings"
)

var _ MetaTable = &PosixMetaTable{}

type PosixMetaTable struct {
	fp       *PlistFile
	parent   *PosixMetaTable
	sub_name string
}

// CreateMetaTable to create or open MetaTable
func CreateMetaTable(path string) (MetaTable, error) {
	if !strings.HasSuffix(path, ".plist") {
		path += ".plist"
	}
	fp, err := CreatePlistFile(path)
	if err != nil {
		return nil, err
	}
	tbl := &PosixMetaTable{
		fp: fp,
	}
	return tbl, nil
}

// CreateMetaTable to open MetaTable
func OpenMetaTable(path string) (MetaTable, error) {
	if !strings.HasSuffix(path, ".plist") {
		path += ".plist"
	}
	fp, err := OpenPlistFile(path)
	if err != nil {
		return nil, err
	}
	tbl := &PosixMetaTable{
		fp: fp,
	}
	return tbl, nil
}

// GetValue return MetaValue according to key
func (tbl *PosixMetaTable) GetValue(v string) MetaValue {
	return tbl.fp.GetValue(v)
}

// SetValue set MetaTable's value,
// plist(darwin, posix): MetaValue ✔ MetaMap ✔ MetaArr ✔
func (tbl *PosixMetaTable) SetValue(v MetaValue) {
	if tbl.parent != nil {
		switch vv := tbl.parent.fp.v.(type) {
		case cfDictionary:
			vv[tbl.sub_name] = v
		case cfArray:

		default:
			tbl.parent.fp.Set(v)
		}
	} else {
		tbl.fp.Set(v)
	}
}

// SafeSetValue set MetaTable's value when key isn't exist,
// plist(darwin, posix): MetaValue ✔ MetaMap ✔ MetaArr ✔
func (tbl *PosixMetaTable) SafeSetValue(v MetaValue) {
	switch vv := v.(type) {
	case MetaMap:
		tbl.fp.Set(cfDictionary{})
		for k, v := range vv {
			tbl.fp.SetByField(k, any2CFValue(v))
		}
	default:
		tbl.SetValue(v)
	}
}

// CreateSubTable create sub element of map or array, but not be saved automatically
func (tbl *PosixMetaTable) CreateSubTable(name string) MetaTable {
	dict := tbl.fp.GetDict()
	if dict.Type() == CF_DICT {
		sub := CFDictionary{}
		dict.Set(name, sub)
		return &PosixMetaTable{
			fp: &PlistFile{
				v:    sub,
				file: Filename(tbl.fp.file) + "." + name,
			},
			parent:   tbl,
			sub_name: name,
		}
	}
	return &PosixMetaTable{}
}

// Write to persist MetaValue in disk.
// note: The regedit (windows) is written when it is created,
// and this method is only valid for plist (darwin, posix).
// The regedit is just an empty method.
func (tbl *PosixMetaTable) Write() error {
	return tbl.fp.Write()
}

// Backup save a copy which could restore MetaValue
func (tbl *PosixMetaTable) Backup() error {
	return tbl.fp.Backup()
}

// Close to free MetaTable memory
func (tbl *PosixMetaTable) Close() {
	tbl.fp.Free()
}
