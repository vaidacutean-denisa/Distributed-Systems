// scrieti un program go care stocheaza in memorie un array cu nr naturale. programul va citi de la tastatura un numar s si va afisa in O(n) indicii elementelor care insumate au ca rezultat s
package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	var arr [10]int = [10]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	nrStr := scanner.Text()
	nr, err := strconv.Atoi(nrStr)

	if err != nil {
		panic(err)
	}
	m := make(map[int]int)

	for idx, val := range arr {
		comp := nr - val

		prevIdx, ok := m[comp]
		if ok {
			fmt.Println(prevIdx, idx)
		}
		m[val] = idx
	}
}
