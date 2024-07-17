package function_call_promise

// 函数调用约定
func A() {
	var a, b int
	b = B(a, b)
}

// 函数A 的栈帧
// 局部变量 a
// 局部变量 b
// 返回值   r
// 参数     b  这里的参数是传递给B 函数的参数，从右到左依次入栈
// 参数     a
// SP of A

// cpu计算的时候会把参数a, b 拷贝到寄存器中计算.
// 在将寄存器上的计算结构拷贝到栈上的返回值空间 r=a+b
// 在go1.17 的版本中是不是这种拷贝，是通过寄存器来传递参数

// int12(1,2,3,3,4,5,6,7,8,9,10,11,12)
// 是用这9个通用寄存器, 寄存器存不下了，会是用栈来继续存储参数
//
// 1  2     3        4  5
// AX,CX,DX,BX,SP,BP,SI,Di
// 6  7  8   9
// R8,R9,R10,R11,R12,R13,R14,R15

// l = 12
// k = 11
// j = 10
// SP os int12

// 这是基础类型的传递方式，如果是结构体传递参数，如果 cpu 的寄存器可以一次 装得下结构体就用 寄存器.
// 如果寄存器存不下，退化为普通的栈传递参数.

// arm64 浮点数不会使用这9个通用寄存器,浮点数是xmm 寄存器 x0-x15
func B(a, b int) (r int) {
	return a + b
}

// B 的函数栈帧
// 返回值 r
// 参数 b
// 参数 a
// SP of B

// go1.17+
// 函数栈
