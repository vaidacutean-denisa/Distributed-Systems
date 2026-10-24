// fie un slice de float32 si un index citit de la tastatura. sa se elimine elementul de la pozitia index, pastrand ordinea celorlalte elemente
package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	slice := []float32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	idxStr := scanner.Text()
	idx, err := strconv.Atoi(idxStr)

	if err != nil {
		panic(err)
	}

	// ... = operator de despachetare a elementelor unui slice, a.i sa poata fi trimise ca argumente individuale lui append
	slice = append(slice[:idx], slice[idx+1:]...)
	fmt.Println(slice)
}
