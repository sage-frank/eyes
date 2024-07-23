package main

func main() {
	x()()
}

func x() (y func()) {
	y = func() {
		println("y")
	}

	return func() {
		println("z")
		y()
	}
}
