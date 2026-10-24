// fie un map de tip string-int care stocheaza perechi de tipul nume-inaltime; sa se afiseze numele persoanelor in ordinea descrescatoare a inaltimilor
package main

import (
	"fmt"
)

func main() {
	m := map[string]int{
		"Andrei": 190,
		"Maria":  165,
		"Ion":    182,
		"Elena":  172,
		"George": 178,
		"Ana":    165,
	}

	// retin un slice cu cheile din map
	var names []string
	for name := range m {
		names = append(names, name)
	}

	// sortez slice-ul folosind criteriul inaltimii asociate
	n := len(names)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if m[names[i]] > m[names[j]] {
				names[i], names[j] = names[j], names[i]
			}
		}
	}

	fmt.Println(names)
}