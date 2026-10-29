// sa se scrie un program in care o goroutina (send) trimite un mesaj unei alte goroutine (receive) printr-un canal
// goroutina receive va returna mesajul, iar main il va afisa
package main

import "fmt"

// canal send-only (trimite la receiver)
func sendMessage(ch1 chan<- string) {
	ch1 <- "message sent"
}

// ch1 e read-only (citeste de la sender); ch2 e send-only (trimite in main)
func receiveMessage(ch1 <-chan string, ch2 chan<- string) {
	message := <-ch1
	ch2 <- message
}

func main() {
	ch1 := make(chan string)
	ch2 := make(chan string)
	go sendMessage(ch1)
	go receiveMessage(ch1, ch2)

	message := <-ch2
	fmt.Println(message)
}
