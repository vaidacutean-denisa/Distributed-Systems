// trei goroutine isi trimit mesaje intre ele prin canale, conform topologiei inel P1 -> P2 -> P3 -> P1
// fiecare dintre aceste evenimente este inregistrat folosind ceasuri logice Lamport. pornind de la sablonul de mai jos, implementati cele 3 metode ale ceasului
// si stocati in main evenimentele ordonate crescator. intre evenimente, contoarele ceasurilor sunt incrementate aleator.
// puteti folosi sort.Slice si functia less

package main

import (
	"fmt"
	"math/rand"
	"sort"
	"sync"
)

type LamportClock struct {
	counter int
}

// incrementeaza contorul
func (l *LamportClock) Tick() {
	l.counter++
}

// ceasul e responsabil de gestionarea ts; procesul trebuie sa se ocupe de comunicarea prin canale
func (l *LamportClock) Send() int {
	l.Tick()
	return l.counter
}

// procesul primeste timestamp-ul mesajului si isi actualizeaza contorul: Cj = max(Cj, ts) + 1 (marcheaza primirea mesajului)
func (l *LamportClock) Receive(ts int) {
	if ts > l.counter {
		l.counter = ts
	}
	l.counter++
}

type Message struct{
	sender,
	ts int
}

type Event struct {
	process int
	internal_counter int
	kind string									// internal, send, recv
}

// returneaza true daca a -> b; ordoneaza mai intai dupa ts, apoi dupa id-ul procesului
func less(a, b Event) bool {
	if a.internal_counter != b.internal_counter {
		return a.internal_counter < b.internal_counter
	}

	return a.process < b.process
}

// functionalitatea unui proces: executa un numar aleatoriu de evenimente + comunica prin canale cu vecinii
func process(id int, clock *LamportClock, in <-chan Message, out chan<- Message, events chan<- Event, wg *sync.WaitGroup) {
	defer wg.Done()
	nr := rand.Intn(3) + 1

	// simulez executia unui nr random de evenimente interne: procesul isi actualizeaza counter-ul
	for i := 0; i < nr; i++ {
		clock.Tick()

		// retin fiecare eveniment in canalul dedicat
		events <- Event{
			id,
			clock.counter,
			"internal",
		}
	}

	// trimite mesaj
	ts := clock.Send()
	events <- Event{
		id,
		ts,
		"send",
	}
	out <- Message{id, ts}

	// primeste mesaj si isi actualizeaza counter-ul
	msg := <-in
	clock.Receive(msg.ts)

	events <- Event{
		id,
		clock.counter,
		"recv",
	}
}

func main() {
	p1 := &LamportClock{}
	p2 := &LamportClock{}
	p3 := &LamportClock{}

	// canale ~ muchii intre nodurile x si y
	ch12 := make(chan Message, 1)
	ch23 := make(chan Message, 1)
	ch31 := make(chan Message, 1)

	events := make(chan Event, 100)
	var wg sync.WaitGroup
	wg.Add(3)

	// pornesc cele trei procese; (id, ceas, canal read, canal send,..)
	go process(1, p1, ch31, ch12, events, &wg)
	go process(2, p2, ch12, ch23, events, &wg)
	go process(3, p3, ch23, ch31, events, &wg)

	wg.Wait()
	close(events)

	var eventList []Event
	for e := range events {
		eventList = append(eventList, e)
	}

	sort.Slice(eventList, func(i, j int) bool {
		return less(eventList[i], eventList[j])
	})

	for _, event := range eventList {
		fmt.Printf("Process P%d, counter = %d, event type = %s\n",
			event.process, event.internal_counter, event.kind)
	}

}

