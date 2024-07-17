package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/bytedance/sonic"
)

func calc(index string, a, b int) int {
	ret := a + b
	fmt.Println(index, a, b, ret)
	return ret
}

func defer1() {
	x := 1
	y := 2

	// 1. a,1,2,3
	//   calc("AA",1,3)
	defer calc("AA", x, calc("a", x, y))

	x = 10

	// 2. b,10,2,12
	// calc("BB", 10, 12)
	// calc("BB", 10, 12) -> calc("AA",1,3)
	defer calc("BB", x, calc("b", x, y))

	y = 10

	// ret
	//  1. a,1,2,3
	// 2. b,10,2,12
	// 3 BB 10,12,22
	// 4 AA 1,3,4
}

//
//func A() {
//	r := deferproc(8, B)   // 通过deferproc负责把要执行的函数信息保存起来,称之为注册
//	if r > 0 {
//		goto ret
//	}
//	// do something
//	runtime.deferreturn  // 返回值前会执行通过 deferproc 注册的函数, 先注册，后调用。
// 多个 defer 会组成一个defer 链表
// 每个 goroutine 执行是都有一个对应的结构体 runtime.g
// g 中有一个 *_defer 的指针指向defer的链表头. defer4 -> defer3->defer2->defer1 新的defer 会加到链表头,执行的时候 也是从头开始执行。
//
//
//	return
//ret:
//	runtime.deferreturn
//}

func deferReturn() int {
	x := 5
	defer func() {
		x++
	}()

	return x
	// 6 同一个x
	// local variable x
	// ret x
	// args
}

func A(a int) {
	fmt.Println(a)
}

func B() {
	x, y := 1, 2
	defer A(x)
	x = x + y
	fmt.Println(x, y)
}

func TrackTime(pre time.Time) time.Duration {
	elapsed := time.Since(pre)
	fmt.Println("elapsed:", elapsed)
	return elapsed
}

func NormalJson() {
	defer TrackTime(time.Now())
	// 定义一个结构体
	type User struct {
		Name   string `json:"name"`
		Age    int    `json:"age"`
		Emails []string
	}

	// 创建一个 User 实例
	user := User{
		Name:   "John Doe",
		Age:    30,
		Emails: []string{"john@example.com", "john.doe@example.com"},
	}

	// 序列化（Marshal）结构体到 JSON 字符串
	jsonBytes, err := json.Marshal(&user)
	if err != nil {
		log.Fatal(err)
	}
	println("JSON:", string(jsonBytes))

	// 反序列化（Unmarshal）JSON 字符串到结构体
	var newUser User
	err = json.Unmarshal(jsonBytes, &newUser)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Unmarshaled User:", newUser)
}

func SonicJson() {
	defer TrackTime(time.Now())
	// 定义一个结构体
	type User struct {
		Name   string `json:"name"`
		Age    int    `json:"age"`
		Emails []string
	}

	// 创建一个 User 实例
	user := User{
		Name:   "John Doe",
		Age:    30,
		Emails: []string{"john@example.com", "john.doe@example.com"},
	}

	// 序列化（Marshal）结构体到 JSON 字符串
	jsonBytes, err := sonic.Marshal(&user)
	if err != nil {
		log.Fatal(err)
	}
	println("JSON:", string(jsonBytes))

	// 反序列化（Unmarshal）JSON 字符串到结构体
	var newUser User
	err = sonic.Unmarshal(jsonBytes, &newUser)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Unmarshaled User:", newUser)
}

func main() {
	// defer1()
	// fmt.Println(deferReturn())
	// B()
	NormalJson()
	SonicJson()
}
