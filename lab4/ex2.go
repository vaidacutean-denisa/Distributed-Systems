// implementati algoritmul de flooding pentru grafuri tare conexe

package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
)

type Message struct {
	Value int
}

func nod(id int, iterations int, in []<-chan Message, out []chan<- Message, isLeader *bool, wg *sync.WaitGroup) {
	defer wg.Done()
	maxId := id

	// executa un numar fix de iteratii de flooding
	for it := 0; it < iterations; it++ {
		// propaga id-ul max pe toate muchiile de iesire
		for _, ch := range out {
			ch <- Message{maxId}
		}

		// actualizez id-ul maxim intalnit
		for _, ch := range in {
			msg := <-ch

			if msg.Value > maxId {
				maxId = msg.Value
			}
		}
	}

	// dupa n-1 iteratii, toate nodurile cunosc id-ul maxim
	if maxId == id {
		*isLeader = true
	}
}

func GraphFlooding(ids []int, graph [][]int) int {
	n := len(graph)

	in := make([][]<-chan Message, n)
	out := make([][]chan<- Message, n)

	// construiesc canale de comunicatie pentru fiecare muchie din graf
	for u := 0; u < n; u++ {
		for _, v := range graph[u] {
			// nodul trimite/citeste cel mult un mesaj de pe fiecare muchie, deci buffer de 1 e suficient
			ch := make(chan Message, 1)
			out[u] = append(out[u], ch)
			in[v] = append(in[v], ch)
		}
	}

	isLeader := make([]bool, n)
	var wg sync.WaitGroup

	// pornesc o goroutina pentru fiecare nod
	for i := 0; i < n; i++ {
		wg.Add(1)
		go nod(ids[i], n - 1, in[i], out[i], &isLeader[i], &wg)
	}
	wg.Wait()

	for i := 0; i < n; i++ {
		if isLeader[i] {
			return ids[i]
		}
	}
	return -1
}

func main() {
	file, err := os.Open("ex2.txt")
	if err != nil { fmt.Println(err) }
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// nr de noduri
	scanner.Scan()
	n, _ := strconv.Atoi(scanner.Text())

	// id-urile nodurilor sunt stocate pe o singura linie
	scanner.Scan()
	idsStr := strings.Fields(scanner.Text())
	ids := make([]int, n)

	for i := 0; i < n; i++ {
		ids[i], _ = strconv.Atoi(idsStr[i])
	}

	// reprezint graful sub forma de liste de adiacenta
	graph := make([][]int, n)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		var u, v int
		fmt.Sscanf(line, "%d %d", &u, &v)

		// noduri numerotate de la 0 la n-1
		u--
		v--
		graph[u] = append(graph[u], v)
	}

	leader := GraphFlooding(ids, graph)
	fmt.Println(leader)
}

// observatia esentiala pe care am omis-o si am pierdut timp incercand sa rezolv deadlock-uri: intr-un graf tare conex, exista drum intre oricare doua noduri
// deci un drum elementar poate avea maxim n-1 muchii => dupa n-1 iteratii, id-ul maxim a ajuns la toate nodurile, proprietate care imi garanteaza un mecanism de oprire

// explicatie extra: la iteratia 1, valoarea maxima pleaca de la nodul de origine si ajunge la toti vecinii directi; la iteratia 2, vecinii dau valoarea mai departe
// prin inductie se demonstreaza ca la iteratia k, valoarea maxima a parcurs o distanta de k muchii fata de sursa;
// intrucat distanta maxima dintre oricare nod si nodul cu idx maxim este <= n-1 muchii, inseamna ca in maxim n-1 iteratii, valoarea se propaga la toate nodurile