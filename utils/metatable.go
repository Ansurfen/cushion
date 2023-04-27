package utils

// MetaTable is an interface to abstract plist(mac) and regedit(windows).
// In other posix os, uniform use of plist as MetaTable interface implement
type MetaTable interface {
	// SetValue set MetaTable's value and have two different rules:
	// regedit(windows): MetaValue ✔ MetaMap ✔ MetaArr x (MetaArr not work);
	// plist(mac, posix): MetaValue ✔ MetaMap ✔ MetaArr ✔
	SetValue(MetaValue)
	// SetValue set MetaTable's value when key isn't exist
	// regedit(windows): MetaValue ✔ MetaMap ✔ MetaArr x (MetaArr not work);
	// plist(mac, posix): MetaValue ✔ MetaMap ✔ MetaArr ✔
	SafeSetValue(MetaValue)
	// GetValue return MetaValue according to key
	GetValue(string) MetaValue
	// CreateSubTable have two different effect:
	// regedit(windows): create sub key and written file depond on its feture.
	// plist(mac, posix): create sub element of map or array, but not be saved automatically
	// comparing regedit. It's required to call Write method save manually.
	CreateSubTable(string) MetaTable
	// Write to persist MetaValue in disk.
	// note: The regedit (windows) is written when it is created,
	// and this method is only valid for plist (mac, posix).
	// The regedit is just an empty method.
	Write() error
	// Backup save a copy which could restore MetaValue
	Backup() error
	// Close to free MetaTable memory
	Close()
}

type (
	MetaValue any
	MetaMap   map[string]any
	MetaArr   []any
)
