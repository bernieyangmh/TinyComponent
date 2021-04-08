package main

import (
	"fmt"
)

func main() {
	diffDeferFunc(1)
}

//defer function is differ with defer func() {function} ()
func diffDeferFunc(i int) (err error) {
	// init
	a1 := &AAA{Name: "aaa", C: &CCC{Age: 18}}
	a2 := AAA{Name: "aaa", C: &CCC{Age: 18}}
	defer fmt.Println(1, i, err, fmt.Sprintf("%+v %+v", a1, a1.C), fmt.Sprintf("%+v %+v", a2, a2.C))
	defer func() { fmt.Println(2, i, err, fmt.Sprintf("%+v %+v", a1, a1.C), fmt.Sprintf("%+v %+v", a2, a2.C)) }()

	// diff
	i = 222
	a1.Name = "a"
	a2.Name = "a"
	a1.C.Age = 16
	a2.C.Age = 16
	err = fmt.Errorf("diff")

	// second
	defer fmt.Println(3, i, err, fmt.Sprintf("%+v %+v", a1, a1.C), fmt.Sprintf("%+v %+v", a2, a2.C))
	defer func() { fmt.Println(4, i, err, fmt.Sprintf("%+v %+v", a1, a1.C), fmt.Sprintf("%+v %+v", a2, a2.C)) }()
	return nil
}

type AAA struct {
	Name string
	C    *CCC
}

type CCC struct {
	Age int
}
