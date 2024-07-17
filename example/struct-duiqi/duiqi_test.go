package struct_duiqi

import (
	"fmt"
	"testing"
	"unsafe"
)

type T struct {
	a int8  // 1byte   8
	b int64 // 8byte   8
	c int32 // 4byte   4
	d int16 // 2byte   4
}

/*
0000 0000  0000 0000 0000 0000 0000 0000
a d  c     b
* /
*/

// 15byte
// 实际是24byte

type T1 struct {
	a int8  // 1byte   0      0
	d int16 // 2byte   2%8=2  2
	c int32 // 4byte   4%8=4  4
	b int64 // 8byte   8%8=0
}

/*
0000 0000  0000 0000 0000 0000 0000 0000
a d  c     b
* /
*/

// 15byte
// 这里调整了一下顺序，实际是16byte就够放下了

func TestDuiQi(t *testing.T) {
	var a T
	fmt.Println(a, unsafe.Sizeof(a))

	var b T1
	fmt.Println(b, unsafe.Sizeof(b))
}
