package cgo

import "C"
import "unsafe"

type (
	CByte   C.char
	CChar   = CByte
	CString *C.char
	CInt    C.int
	CUInt   C.uint
	CInt8   C.schar
	CUInt8  C.uchar
	CInt16  C.short
	CUInt16 C.ushort
	CLong   C.long
	CULong  C.ulong
	CInt64  C.longlong
	CUInt64 C.ulonglong
	CFloat  C.float
	CDouble C.double
	CSize_t = C.size_t
	Ptr     = unsafe.Pointer
	Voidptr = uintptr
)

func GStr(str CString) string {
	return C.GoString(str)
}

func CStr(str string) CString {
	return C.CString(str)
}

func CastPtr[T any](ptr Ptr) *T {
	return (*T)(ptr)
}

func CastVoidPtr[T any](a *T) Voidptr {
	return Voidptr(Ptr(a))
}

func Cast[T any](ptr any) T {
	return ptr.(T)
}

func CMalloc(size CSize_t) Ptr {
	return C.malloc(size)
}

func Nullptr() Ptr {
	return C.NULL
}
