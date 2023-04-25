//go:build !windows
// +build !windows

package utils

import (
	"strings"
)

type PosixMetaTable struct {
	fp       *PlistFile
	parent   *PosixMetaTable
	sub_name string
}

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

func (tbl *PosixMetaTable) GetValue(v string) MetaValue {
	return tbl.fp.GetValue(v)
}

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

func (tbl *PosixMetaTable) Write() error {
	return tbl.fp.Write()
}

func (tbl *PosixMetaTable) Backup() error {
	return tbl.fp.Backup()
}

func (tbl *PosixMetaTable) Close() {
	tbl.fp.Free()
}
