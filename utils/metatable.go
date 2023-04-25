package utils

type MetaTable interface {
	SetValue(MetaValue)
	SafeSetValue(MetaValue)
	GetValue(string) MetaValue
	CreateSubTable(string) MetaTable
	Write() error
	Backup() error
}

type (
	MetaValue any
	MetaMap   map[string]any
	MetaArr   []any
)
