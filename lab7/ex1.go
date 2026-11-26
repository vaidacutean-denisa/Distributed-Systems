// sa se implementeze algoritmul ByzFlood pentru un sistem distribuit sincron cu n noduri si cel mult s defecte bizantine (n > 4s), reprezentat printr-un graf complet
package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
)

type Node struct {
	id int
	x int								// starea curenta xi(t)
	s int
	n int								// nr de noduri din sistem
	isByzantine bool					// comportamentul la iteratia curenta
	inbox chan Message
}

// mesaj transmis intre noduri in runda 1 sau in runda 2
type Message struct {
	SenderID int
	Round int					// =tipul de mesaj -> nodul trebuie sa poata distinge intre starea votata de alt nod (round=1) sau decizia tie trimisa de lider (round=2)
	Value int					// 0 sau 1 ~ votul
}

func firstRound(idx int, nodes[]* Node) (int, int){
	node := nodes[idx]

	// daca nodul este bizantin, trimite valoarea opusa
	val := node.x
	if node.isByzantine {
		val = 1 - val
	}

	// etapa de broadcast: nodul curent trimite mesaj fiecarui vecin (graf complet)
	for j := 0; j < node.n; j++ {
		msg := Message{
			node.id,
			1,
			val,
		}
		nodes[j].inbox <- msg
	}

	// etapa de receive: nodul calculeaza majority(t) = maj(U), cu U = multimea mesajelor primite de la vecini; retine si frecventa acestei valori
	var frv_ones int = 0
	for i := 0; i < node.n; i++ {
		msg := <-node.inbox
		if msg.Round == 1 && msg.Value == 1 {
			frv_ones++
		}
	}

	var maj int = 1
	var frv int = frv_ones

	if frv_ones < len(nodes) / 2 + 1{
		maj = 0
		frv = len(nodes) - frv_ones
	}
	return maj, frv
}

func secondRound(s int, t int, idx int, nodes[]* Node, maj int, frv int) {
	node := nodes[idx]

	tie := maj
	// daca nodul curent are idx = t, este lider si trimite propria valoare maj; altfel primeste valoarea trimisa de lider
	if t == idx {
		for i := 0; i < len(nodes); i++ {
			if i == idx {
				continue
			}

			// daca nodul curent e bizantin, trimite valori tie random (0/1)
			if node.isByzantine {
				tie = rand.Intn(2)
			}

			msg := Message{
				idx,
				2,
				tie,
			}
			nodes[i].inbox <- msg
		}
		// liderul isi actualizeaza starea
		node.x = tie
	} else {
		// daca nu e lider, nodul isi actualizeaza starea in functie de frecventa propriei majoritati si valoarea primita de la lider
		msg := <-node.inbox
		if msg.Round == 2 {
			recv_tie := msg.Value

			node.x = recv_tie
			if frv > len(nodes)/2+s {
				node.x = maj
			}
		}
	}
}

func byzFlood(nodes []*Node, s int) {
	n := len(nodes)
	maj := make([]int, n)
	frv := make([]int, n)

	for t := 0; t < s + 1; t++ {
		var wg sync.WaitGroup

		// runda 1 (astept sa termine toate nodurile runda 1 inainte de a trece la runda 2; altfel ar fi fost posibil sa se amestece mesajele dintre runde in inbox)
		for i := 0; i < n; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				maj[i], frv[i] = firstRound(i, nodes)
			}(i)
		}
		wg.Wait()

		// runda 2
		for i := 0; i < n; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				secondRound(s, t, i, nodes, maj[i], frv[i])
			}(i)
		}
		wg.Wait()
	}

	// afisez rezultatul
	for _, nd := range nodes {
		if !nd.isByzantine {
			fmt.Println("Majority: ", nd.x)
			// obs: am testat fara break ca sa verific ca toate nodurile corecte au aceeasi valoare
			break
		}
	}
}

func main() {
	file, err := os.Open("ex1.txt")
	if err != nil { fmt.Println(err) }
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	values := strings.Fields(scanner.Text())

	n, _ := strconv.Atoi(values[0])
	s, _ := strconv.Atoi(values[1])

	nodes := make([]*Node, n)
	for i := 0; i < n; i++ {
		scanner.Scan()
		vote, _ := strconv.Atoi(scanner.Text())
		nodes[i] = &Node {
			id: i,
			x: vote,
			s: s,
			n: n,
			isByzantine: false,
			inbox: make(chan Message, n),
		}
	}

	// generez un numar random de defecte (cel mult s); generez o permutare aleatorie a indicilor -> primii f vor fi bizantini
	f := rand.Intn(s + 1)
	perm := rand.Perm(n)

	for i := 0; i < f; i++ {
		nodes[perm[i]].isByzantine = true
	}
	byzFlood(nodes, s)
}