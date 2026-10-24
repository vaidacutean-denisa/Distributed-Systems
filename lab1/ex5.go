// fie un slice de string-uri. sa se afiseze elementele acestuia in ordine inversa
package main

import "fmt"

func main() {
	slice := []string{"Laboratorul", "de", "Sisteme", "Distribuite"}

	n := len(slice)
	// pornesc cu i de la 0 si j de la n-1; fac swap intre elementele slice-ului
	for i, j := 0, n - 1; i < j; i, j = i + 1, j - 1 {
		slice[i], slice[j] = slice[j], slice[i]
	}

	fmt.Println(slice)
}
