// implementati algoritmul de flooding pentru un inel unidirectional

package main

import (
	"fmt"
	"sync"
)

type Msg struct {
	Value int
	Done bool
}

// realizeaza operatiile interne nodului
func node(id int, in <-chan Msg, out chan<- Msg, isLeader *bool, wg *sync.WaitGroup) {
	defer wg.Done()

	// propaga propriul id in inel
	out <- Msg{Value: id, Done: false}

	for msg := range in {
		// daca a gasit liderul, propaga mesajul si iese
		if msg.Done {
			out <- msg
			break
		}

		// valoarea maxima s-a intors la nodul asociat -> devine lider si propaga mesajul de oprire
		if msg.Value == id {
			*isLeader = true
			out <- Msg{Value: id, Done: true}
			break
		} else if msg.Value > id {						// id-ul primit e mai mare -> retrimite mesajul; id mai mic => null
			out <- msg
		}
	}
}

// returneaza id-ul liderului
func RingFlooding(ids []int) int {
	n := len(ids)

	// construiesc canale de comunicatie pentru fiecare muchie din inel
	channels := make([]chan Msg, n)
	for i := 0; i < n; i++ {
		// canale bufferizate; fiecare nod incearca mai intai sa trimita propriul id, deci nimeni nu face <-chan => deadlock
		channels[i] = make(chan Msg, 1)
	}

	isLeader := make([]bool, n)
	wg := sync.WaitGroup{}

	// nodul i citeste de la predecesorul sau (i-1+n mod n) si trimite la succesor
	for i := 0; i < n; i++ {
		wg.Add(1)
		in := channels[(i - 1 + n) % n]
		out := channels[i]

		go node(ids[i], in, out, &isLeader[i], &wg)
	}
	// asteapta toate goroutinele ~ contorul ajunge la 0
	wg.Wait()

	for i := 0; i < n; i++ {
		if isLeader[i] == true {
			return ids[i]
		}
	}
	return -1
}

func main() {
	ids := []int{7, 3, 1, 9, 5, 4, 2}
	fmt.Println(RingFlooding(ids))
}
