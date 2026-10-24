// angajatii unei firme sunt stocati intr-un dictionar map[string]int cu perechi de forma nume-id_departament
// creati un nou map[int][]string care sa stocheze perechi de tipul id_dept-slice cu numele angajatilor din departament
package main

import (
	"fmt"
	"strings"
)

func main() {
	 emp := map[string]int {
		 "Ana Popescu": 1,
		 "Mihai Ionescu": 2,
		 "Bogdan Radu": 1,
		 "Elena Dumitrescu": 3,
		 "Cristian Vasile": 2,
		 "Ioana Enache": 1,
	 }

    dept := make(map[int][]string)
	for name, id := range emp {
		dept[id] = append(dept[id], name)				// daca id-ul nu exista inca, dept[id] are valoarea default pentru tipul []string
	}

	for id, name := range dept {
		fmt.Printf("Departamentul %d: %s\n", id, strings.Join(name, ", "))
	}
}