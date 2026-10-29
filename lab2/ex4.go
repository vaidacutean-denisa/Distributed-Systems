// sa se scrie un program in care o goroutina sa notifice functia main atunci cand procesarea sa se sfarseste, fara a folosi un wait group
package main

import (
	"fmt"
)

func process(ch chan string) {
	defer close(ch)									// inchide canalul dupa ce isi finalizeaza executia
	fmt.Println("Currently in process..")
}

func main() {
	ch := make(chan string)
	go process(ch)

	<-ch											// main se blocheaza pana cand process trimite ceva pe canal sau il inchide
	fmt.Println("Now in main")
}