package main

import (
	"fmt"
)

func main() {
	ChildenShareParentVariables()
	fmt.Println("\nafter fix\n")
	Fix_ChildenShareParentVariables()
}

// children G share i，parent G's variables
// when children G start， maybe the i already has changed
func ChildenShareParentVariables() {
	for i := 0; i < 100; i++ { // write
		go func() { /* Create a new goroutine */
			go func(i int) {
				fmt.Println(i)
			}(i)
		}()
	}
}

func Fix_ChildenShareParentVariables() {
	for i := 0; i < 100; i++ { // write
		tmp := i
		go func(i int) { /* Create a new goroutine */
			go func() {
				fmt.Println(i)
			}()
		}(tmp)
	}
}
