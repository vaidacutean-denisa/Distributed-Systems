// sa se scrie un program in care goroutina main primeste doua mesaje pe acelasi canal de comunicare
package main

import "fmt"

func sendFirstMsg(ch chan string) {
	ch <- "first message sent"
}

// obs: nu am garantia ca goroutina sendFirst va trimite prima mesajul
func sendSecondMsg(ch chan string) {
	ch <- "second message sent"
}

func main() {
	ch := make(chan string)
	go sendFirstMsg(ch)
	go sendSecondMsg(ch)

	for i := 0; i < 2; i++ {
		fmt.Println(<-ch)
	}
}
