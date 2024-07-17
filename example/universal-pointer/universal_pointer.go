package main

import "fmt"

func defineVar() {
	s := 123
	fmt.Printf("%p, %d\n", &s, *(&s)) // 0xc000096068, 123

	var j int
	fmt.Printf("%p, %d\n", &j, *(&j)) // 0xc00000a100, 0 //定义变量后就分配了内存地址，但是没有赋值
}

func saveAddr() {
	var s string = "abc"

	var p *string = &s

	p1 := &s

	fmt.Println(p, p1) // 0xc000024070 0xc000024070  p 和 p1 的值是同一个，指向同一个地址

	fmt.Println(*p) // abc 通过 *取出指向地址的值

	*p = "world"
	fmt.Println(*p1) // p1 和p 指向的同一个地址空间，所以 *p改了值，*p1 的值也跟着改
}

func DiffPointer() {
	a := 123
	b := a
	// b = a 表示将 a 的值拷贝一份给 b, 所以这两个变量的地址是不一样的
	fmt.Printf("%p %p\n", &a, &b) // 0xc00000a110 0xc00000a118

	c := 456
	p1 := &c
	p2 := p1
	// 因为将 p1 赋值给了 p2, 所以这两个指针变量的值是一样的, 因为存储的都是变量 c 的地址
	// 但我们说了, Go 变量的传递都是值传递, 也就是要拷贝一份出来
	// 所以 p1 和 p2 存储的值一样, 表示它们存储的地址是一样的, 修改 *p1 会影响 *p2, 修改 *p2 会影响 *p1, 因为都指向同一份内存
	// 但是这两个指针变量本身的地址是不一样的, 指针变量也是有地址的
	fmt.Printf("%p %p\n", &p1, &p2) // 0xc000064030 0xc000064038
	// 此时修改 p1, 不会影响 p2; 修改 p2 不会影响 p1, 因为它们是两个不同的变量, 只是存储的值(地址)一样, 而其本身的地址不一样
}

func NullPointer() {
	fmt.Println("------------")
	var p1 *int
	var p2 *int
	fmt.Printf("%p, %p, %t\n", p1, p2, p1 == p2) // panic: runtime error: invalid memory address or nil pointer dereference
}

func NewPointer() {
	p := new(int)
	fmt.Printf("%p,%d", p, *p) // 0xc000096068,0 分配了指针地址，也分配了指针指向地址 默认值初始值为0
	*p = 123
	fmt.Println(p, *p)

	// 上面就等价于:
	var i int
	p2 := &i
	fmt.Println(*p2, i) // 0

	// 只不过此时我们知道 p2 这个指针变量指向谁, 修改其中一个都会影响另一个, 因为都是同一份内存
	i = 1
	fmt.Println(*p2, i) // 1 1
	*p2 = 2
	fmt.Println(*p2, i) // 2 2
}

func main() {
	// saveAddr()
	// DiffPointer()
	// NullPointer()
	NewPointer()
}
