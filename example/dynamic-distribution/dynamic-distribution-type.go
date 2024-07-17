package main

import "fmt"

type empty interface{}

// iface
//
//type iface struct {
//	tab     *itab
//	pointer unsafe.Pointer
//}
//
//type itab struct {
//	inter *interfacetype
//	_type *_type
//	hash  uint32
//	_     [4]byte
//	fun   [1]uintptr
//}

type MyType int

type Integer interface {
	~int8 | ~int16 | ~int | ~int32 | ~int64
}

func Sum[V Integer](values []V) V {
	var ret V
	for _, v := range values {
		ret += v
	}
	return ret
}

func Print[T fmt.Stringer](v T) {
	fmt.Println(v)
}

func main() {
	Sum([]int{1, 2, 3})
}
