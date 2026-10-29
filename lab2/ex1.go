// definiti un program go care sa transfere cu succes un mesaj de la o goroutina diferita de main() la goroutina main() prin intermediul unui channel
package main

import (
	"fmt"
)

func sendMsg(ch chan string) {
	ch <- "goroutina sendMsg"
}

func main() {
	ch := make(chan string)
	go sendMsg(ch)
	fmt.Println(<-ch)
}