// fie un sistem distribuit format din doi roboti identici, porniti de catre un observer (main). fiecare robot parcurge un traseu si trimite periodic coord sale 2D catre celalalt robot, in timp ce se deplaseaza.
// cand un robot isi termina de parcurs traseul, il va notifica pe celalalt ca s-a oprit. acesta trebuie sa citeasca toate pozitiile primite pana la oprire, iar apoi ambii vor raporta observer-ului ca si-au terminat sarcina
// sa se simuleze comportamentul robotilor

package main

import (
	"fmt"
	"sync"
	"time"
)

type Coordinates struct {
	X int
	Y int
}

// logica asociata unui robot: citeste datele primite de la celalalt robot si, in paralel, trimite propriile coordonate in timp ce isi parcurge traseul
func processRobot(id int, path []Coordinates, send chan<- Coordinates, receive <-chan Coordinates, done chan<- bool) {
	var wg sync.WaitGroup
	wg.Add(1)

	// goroutina anonima care asculta pe canalul read-only, in paralel; cu wg astept sa isi termine executia, inainte sa ies din functia principala
	go func() {
		defer wg.Done()
		for pos := range receive {
			fmt.Printf("Robot %d received: (%d, %d)\n\n", id, pos.X, pos.Y)
		}
		fmt.Printf("Robot %d stopped listening on the channel\n\n", id)
	}()

	// robotul parcurge traseul in acelasi timp
	for _, coord := range path {
		time.Sleep(100 * time.Millisecond)
		send <- coord
		fmt.Printf("Robot %d sent: (%d, %d)\n", id, coord.X, coord.Y)
	}

	// anunta celalalt robot ca a finalizat parcurgerea traseului prin inchiderea canalului de send
	close(send)
	wg.Wait()
	done <- true
}

func main() {
	// ch(xy): send-only pentru x si read-only pentru y (ca o muchie unidirectionala intre nodurile 1 si 2, unde nod=robot)
	ch12 := make(chan Coordinates)
	ch21 := make(chan Coordinates)

	// main asteapta fiecare robot
	done1 := make(chan bool)
	done2 := make(chan bool)

	path1 := []Coordinates{
		{X: 1, Y: 4},
		{X: 2, Y: 3},
	}

	path2 := []Coordinates{
		{X: 7, Y: 8},
		{X: 9, Y: 7},
	}

	go processRobot(1, path1, ch12, ch21, done1)				// mai explicit, argumentele sunt (id, path, **send, receive**, done)
	go processRobot(2, path2, ch21, ch12, done2)
	<-done1
	<-done2

	fmt.Println("Observer: both robots finished their task.\n")
}