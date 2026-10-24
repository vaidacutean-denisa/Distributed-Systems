// sa se citeasca dintr-un fisier un slice de string-uri formate numai din litere mici. sa se afiseze frecventa fiecarui cuvant in slice
package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	file, err := os.Open("ex7.txt")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()											// defer amana inchiderea fisierului pana la terminarea functiei curente
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)							    // toate cuvintele sunt pe aceeasi linie

	m := make(map[string]int)
	for scanner.Scan() {
		str := scanner.Text()
		m[str]++;
	}

	err = scanner.Err()
	if err != nil {
		fmt.Println(err)
	}

	for word, count := range m {
		fmt.Println(word, count)
	}
}