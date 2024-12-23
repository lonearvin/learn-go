package main

func main() {
	test()
}

func test() {
	defer func() {
		if err := recover(); err != nil {
			println(err.(string))
		}
	}()
	panic("panic error")
}
