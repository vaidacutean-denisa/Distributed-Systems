// fie un array in memorie. sa se sorteze crescator si sa se afiseze
package main

import (
	"fmt"
)

func main() {
	var intArr [5]int32 = [5]int32{7, 3, 1, 9, 4}

	for i := 0; i < 4; i++ {
		for j := 0; j < 4-i; j++ {
			if intArr[j] > intArr[j+1] {
				intArr[j], intArr[j+1] = intArr[j+1], intArr[j]
			}
		}
	}
	fmt.Println(intArr)
}
