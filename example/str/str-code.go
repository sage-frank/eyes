package main

import (
	"fmt"
	"unsafe"
)

// 字符串的结构应该是这样的
//type String struct {
//	p   *unsafe.Pointer // 8
//	len int             // 8
//}

func main1() {
	fmt.Println(unsafe.Sizeof("abc"))  // 16
	fmt.Println(unsafe.Sizeof("中文字符")) // 16  字符串的长度都是16，显然这个变量没有实际存字符串的内容,
	s := "你好么?OK"
	fmt.Println(s[0:3], s[3:6], s[6:], []byte(s))
	// 中英文混合肯定不行utf-8 的编码，英文字符是1个字节，中文是3个字节
	toRune(s)

	utf8Code(s)

	ToUtf8([]byte{228, 189, 160, 229, 165, 189, 228, 185, 136, 63, 79, 75},
		[]rune{20320, 22909, 20040, 63, 79, 75})
}

func toRune(s string) {
	s1 := []rune(s)
	fmt.Println(s1)
	fmt.Println(string(s1[0]), string(s1[1]), string(s1[2]))
}

func utf8Code(s string) {
	fmt.Println([]byte(s)) //[228 189 160 229 165 189 228 185 136 63 79 75]
	fmt.Println([]rune(s)) //[20320 22909 20040 63 79 75]
}

func ToUtf8(b1 []byte, r1 []rune) {
	fmt.Println(string(b1))
	fmt.Println(string(r1))
}

// 字符串和切片
// runtime/string.go
/*
type stringStruct struct {
	str unsafe.Pointer
	len int
}

// runtime/slice.go
type slice struct {
    array unsafe.Pointer
    len   int
    cap   int
}
*/

func int8Toint16() {
	s1 := []int8{1, 2, 3, 4}
	//  int8 不能转为int16 会报错
	// s2 := []int16(s1) //  cannot convert s1 (variable of type []int8) to type []int16
	// fmt.Println(s2)
	s2 := *(*[]int16)(unsafe.Pointer(&s1))
	fmt.Println(s2) // [513 1027 0 0]
	// 1,2
	// 3,4
	fmt.Println(2<<8 + 1) // 513
	fmt.Println(4<<8 + 3) // 1027
}

func strShareSlice() {
	fmt.Println("------")
	str := "abc"
	slice := *(*[]byte)(unsafe.Pointer(&str))
	fmt.Println(slice)      // [97 98 99]
	fmt.Println(cap(slice)) // 10036576
}

func StringToBytes(s string) []byte {
	// 既然字符串转切片，会丢失容量
	// 那么加上去就好了，做法也很简单
	// 新建一个结构体，将容量（等于长度）加进去
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{s, len(s)},
	))
}

func BytesToString(b []byte) string {
	// 切片转字符串就简单了，直接转即可
	// 转的过程中，切片的 Cap 字段会丢弃
	return *(*string)(unsafe.Pointer(&b))
}

// AssertShareArray 证明是共享数组了
func AssertShareArray() {
	slice := []byte{97, 98, 99}
	str := *(*string)(unsafe.Pointer(&slice))
	fmt.Println(str) // abc
	slice[0] = 'A'
	fmt.Println(str) // Abc
}

func main() {
	// int8Toint16()
	// strShareSlice()
	// slice1 := StringToBytes("abc")
	// fmt.Println(slice1, len(slice1), cap(slice1))
	// fmt.Println(BytesToString([]byte{97, 98, 99}))

	AssertShareArray()
}
