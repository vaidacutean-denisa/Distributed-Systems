// un client cu ceasul decalat cere de la un server timpul corect printr-o retea cu intarzieri simulate cu time.Sleep
// pe langa θ, vom folosi si durata round-trip-uli, fara timpul de procesare al serverului: S = (T4 - T1) - (T3 - T2)
// a) sa se completeze functia estimate folosind algoritmul lui Cristian; sa se simuleze in main o cerere si un raspuns folosind measure; sa se verifice daca offset-ul initial a scazut
// b) sa se repete simularea de 10 ori (cu acelasi offset initial) si calculati S si θ minime. comparati rezultatele cu cele obtinute la a)

package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Clock struct {
	offset time.Duration
}

func (c *Clock) Now() time.Time {
	return time.Now().Add(c.offset)
}

// algoritmul lui Cristian
func estimate (t1 , t2 , t3 , t4 time.Time) (theta, delta time.Duration) {
	theta = ((t2.Sub(t1)) + (t3.Sub(t4))) / 2				// (t2 - t1 + t3 - t4) / 2; time nu suporta op -
	delta = (t4.Sub(t1)) - (t3.Sub(t2))						// S = t4 - t1 - (t3 - t2)

	return theta, delta
}

func measure (c_clock , s_clock *Clock) (theta, delta time.Duration) {
	t1 := c_clock.Now()
	time.Sleep (3* time.Millisecond + time. Duration(rand.Intn(10))*time.Millisecond)

	t2 := s_clock.Now()
	time.Sleep(5* time.Millisecond + time.Duration(rand.Intn(20))*time.Millisecond)

	t3 := s_clock.Now()
	time.Sleep(3* time.Millisecond + time.Duration(rand.Intn(10))*time.Millisecond)

	t4 := c_clock.Now()
	return estimate(t1 , t2 , t3 , t4)
}

func main() {
	rand.Seed(time.Now().UnixNano())

	// a) o singura simulare
	initialOffset := 100 * time.Millisecond

	client := &Clock{initialOffset}
	server := &Clock{0}
	fmt.Println("Initial offset: ", initialOffset)

	theta, delta := measure(client, server)
	fmt.Println("Theta: ", theta)
	fmt.Println("Delta: ", delta)

	// corectez ceasul clientului
	client.offset += theta

	fmt.Println("Updated offset: ", client.offset)
	fmt.Println("The offset has decreased: ", client.offset.Abs() < initialOffset.Abs())

	// b) 10 simulari; initializez min theta si delta cu valorile maxime pentru int64
	minTheta := time.Duration(1 << 63 - 1)
	minDelta := time.Duration(1 << 63 - 1)

	for i := 1; i <= 10; i++ {
		client.offset = initialOffset
		theta, delta := measure(client, server)

		fmt.Printf("Iteration %d: theta = %v, delta = %v\n", i, theta, delta)

		if theta < minTheta {minTheta = theta}
		if delta < minDelta {minDelta = delta}
	}
	fmt.Println("Min theta: ", minTheta)
	fmt.Println("Min delta: ", minDelta)
}

// observatii: in prima simulare (punctul a) am obtinut delta de ~7.4ms, iar in urmatoarele 10 simulari (punctul b), delta minim a fost ~10.3ms

// pentru a)
// clientul are offset initial +100ms; serverul are 0
// offset = server - client => obtin offset aprox egal cu -100.529ms < 0, deci clientul este inaintea serverului; adun offset-ul la ceasul clientului, obtin noul offset de aprox -0.529ms
// concluzie: clientul a trecut de la decalaj de 100ms la decalaj de aproximativ 0.53ms
