package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
)

// bfs pornit din fiecare nod; returneaza -1 daca graful nu este tare conex
func diameter(adj [][]int) int {
	n := len(adj)
	diam := 0

	for src := 0; src < n; src++ {
		dist := make([]int, n)
		for i := range dist {
			dist[i] = -1
		}

		dist[src] = 0
		queue := []int{src}

		for len(queue) > 0 {
			u := queue[0]
			queue = queue[1:]

			// actualizeaza distanta de la src pana la vecinii nevizitati
			for _, v := range adj[u] {
				if dist[v] == -1 {
					dist[v] = dist[u] + 1
					queue = append(queue, v)
				}
			}
		}
		for _, d := range dist {
			// nod inaccesibil din src
			if d == -1 {
				return -1
			}
			// actualizeaza diametrul grafului
			if d > diam {
				diam = d
			}
		}
	}
	return diam
}

type Info struct {
	ID int
	Value int
}

type Memory map[Info]struct{}

type Node struct {
	ID int
	Neighbours []int
	Mem Memory
	IsActive bool

	Out []chan<- Memory
	In []<-chan Memory
}

func Crash(nodes []Node) bool {
	crash := rand.Intn(2) == 0
	if crash {
		var activeNodes []int
		for idx, node := range nodes {
			if node.IsActive {
				activeNodes = append(activeNodes, idx)
			}
		}
		if len(activeNodes) == 0 {
			return false
		}
		randIdx := rand.Intn(len(activeNodes))
		nodes[activeNodes[randIdx]].IsActive = false
		return true
	}
	return false
}

// nodul curent: 1) trimite memoria sa catre toti vecinii; 2) face reuniunea dintre memoria proprie si cea primita de la vecini
func ProcessNode(node int, nodes []Node, adj [][]int, wg *sync.WaitGroup) {
	defer wg.Done()

	// etapa 1: nodul trimite memoria sa pe canalele de iesire sau nil daca e mort (valoarea implicita pentru map e nil)
	var content Memory
	if nodes[node].IsActive {
		content = nodes[node].Mem
	}
	for _, ch := range nodes[node].Out {
		ch <- content
	}

	// daca nodul e inactiv, eliberez canalele de input si ies
	if !nodes[node].IsActive {
		for _, ch := range nodes[node].In {
			<-ch
		}
		return
	}

	// etapa 2: creez un memory nou pentru nod; nu actualizez mem dinamic pentru ca alte noduri o pot citi in acelasi timp
	newMemory := make(Memory)

	// copiez memoria veche
	for info := range nodes[node].Mem {
		newMemory[info] = struct{}{}
	}
	// nodul citeste datele primite de la vecini si le adauga la memoria noua (face union automat -> suprascrie valoarea deja existenta)
	for _, ch := range nodes[node].In {
		msg := <-ch
		if msg != nil {
			for info := range msg {
				newMemory[info] = struct{}{}
			}
		}
	}
	nodes[node].Mem = newMemory
}

func FloodSet(nodes []Node, adj [][]int, s int, limit int, diam int) {
	var crashedCount int = 0

	for t := 0; t <= limit; t++ {
		// daca nu am depasit totalul de noduri inactive, decid random daca modific statusul unuia dintre cele ramase
		if crashedCount < s {
			crashed := Crash(nodes)
			if crashed {
				crashedCount++
			}
		}

		var wg sync.WaitGroup
		for i := 0; i < len(nodes); i++ {
			wg.Add(1)
			go ProcessNode(i, nodes, adj, &wg)
		}
		wg.Wait()
	}

	fmt.Printf("diameter: %d\n", diam)
	fmt.Printf("iterations: %d\n", limit)

	// caut primul nod activ pentru a afisa memoria
	activeIdx := -1
	for idx, node := range nodes {
		if node.IsActive {
			activeIdx = idx
			break
		}
	}
	if activeIdx != -1 {
		var result []Info
		for info := range nodes[activeIdx].Mem {
			result = append(result, info)
		}
		sort.Slice(result, func(i, j int) bool {
			return result[i].ID < result[j].ID
		})

		fmt.Print("Memory: ")
		for _, info := range result {
			fmt.Printf("(%d, %d) ", info.ID, info.Value)
		}
		fmt.Println()
	}
}

func main() {
	file, err := os.Open("ex1.txt")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var n, s int

	if scanner.Scan() {
		line := scanner.Text()
		fmt.Sscanf(line, "%d %d", &n, &s)
	}

	// citesc tokenii si creez nodurile
	nodes := make([]Node, n)
	if scanner.Scan() {
		tokens := strings.Fields(scanner.Text())
		for i := 0; i < n; i++ {
			val, _ := strconv.Atoi(tokens[i])

			info := Info{i, val}
			mem := make(Memory)
			mem[info] = struct{}{}

			nodes[i] = Node{
				ID: i,
				Mem: mem,
				IsActive: true,
			}
		}
	}

	// creez listele de adiacenta
	adj := make([][]int, n)
	for i := range adj {
		if scanner.Scan() {
			neighbours := strings.Fields(scanner.Text())

			for _, neighbour := range neighbours {
				nr, _ := strconv.Atoi(neighbour)
				adj[i] = append(adj[i], nr)
			}
			nodes[i].Neighbours = adj[i]
		}
	}

	// fiecare nod are asociata cate o lista de canale in/out; le initializez
	for i := 0; i < n; i++ {
		nodes[i].Out = make([]chan<- Memory, 0)
		nodes[i].In = make([]<-chan Memory, 0)
	}

	// mapez muchiile pe canale
	diam := diameter(adj)
	if diam == -1 {
		fmt.Println("the graph is not strongly connected")
		return
	}
	limit := (s + 1) * diam

	for u := 0; u < n; u++ {
		for _, v := range adj[u] {
			ch := make(chan Memory, limit + 1)

			nodes[u].Out = append(nodes[u].Out, ch)
			nodes[v].In = append(nodes[v].In, ch)
		}
	}
	FloodSet(nodes, adj, s, limit, diam)
}
