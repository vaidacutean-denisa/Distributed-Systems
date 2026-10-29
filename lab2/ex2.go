// sa se transfere un mesaj de la goroutina main la alta, diferita de main, prin intermediul unui channel
package main

import (
	"fmt"
	"sync"
)

func receiveMsg(ch chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println(<-ch)
}

func main() {
	ch := make(chan string)
	var wg sync.WaitGroup

	wg.Add(1)						// mecanismul de sincronizare previne ca main sa-si finalizeze executia inainte ca receiveMsg sa afiseze mesajul
	go receiveMsg(ch, &wg)				// pornesc o goroutina concurenta, apoi trimit mesajul; altfel primesc deadlock (presupunand un canal nebufferizat)
	ch <- "goroutina main"

	wg.Wait()
}

