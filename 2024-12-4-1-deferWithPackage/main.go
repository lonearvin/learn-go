package main

import "fmt"

type Test struct {
	name string
}

func (t *Test) Close() {
	fmt.Println(t.name)
}

func main() {

	test := []Test{{"a"}, {"b"}, {"c"}}

	for _, v := range test {
		defer v.Close()
	}

}
