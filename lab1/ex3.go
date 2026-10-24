// fie un sir de numere. sa se afiseze un slice cu numerele pare din acest sir
package main

import (
	"fmt"
)

func main() {
	var arr [10]int32 = [10]int32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	slice := []int32{}
	for i := range arr {
		if arr[i] % 2 == 0 {
			slice = append(slice, arr[i])
		}
	}
	fmt.Println(slice)
}
