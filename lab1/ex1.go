// se citeste un numar natural de la tastatura; sa se verifice daca este prim
package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	// ia ca parametru un scanner si returneaza un bool
	scanner.Scan()

	// conversia in string la ce am citit mai sus
	nrStr := scanner.Text()
	nr, err := strconv.Atoi(nrStr)

	if err != nil {
		panic(err)
	}

	if nr < 2 {
		fmt.Printf("Numarul nu este prim\n")
		return
	}

	for i := 2; i * i <= nr; i++ {
		if nr % i == 0 {
			fmt.Printf("Numarul %v nu este prim", nr)
			fmt.Printf("\n")
			return
		}
	}
	fmt.Printf("Numarul %v este prim", nr)
	fmt.Printf("\n")
}
