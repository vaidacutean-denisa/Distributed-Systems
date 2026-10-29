// fie un sistem distribuit format din doi roboti identici. ei sunt porniti de catre un observer, reprezentat de functia main. cei doi roboti pornesc in acelasi timp si se intrec in linie dreapta pana la un finish
// sa se scrie un program care simuleaza acest sistem, unde functia main va afisa care robot a ajuns primul la finish. hint: folositi select
package main

import (
	"fmt"
	"time"
)

func robot(name string, finish chan<- string) {
	time.Sleep(2 * time.Second)
	finish <- name
}

func main() {
	ch1 := make(chan string)
	ch2 := make(chan string)

	go robot("Robot 1", ch1)
	go robot("Robot 2", ch2)

	select {										// select se blocheaza pana cand cel putin unul dintre canale primeste un mesaj
	case winner := <-ch1:						    // primul care trimite mesaj castiga, desi daca trimit simultan, select alege random un bloc de cod pe care sa-l execute
		fmt.Println(winner)
	case winner := <-ch2:
		fmt.Println(winner)
	}
}