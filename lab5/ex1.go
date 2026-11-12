// sa se implementeze algoritmul de averaging (problema rendezvous)

package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"sync"
)

type Coord struct {
	X float64
	Y float64
}

// genereaza G=(V,E) neorientat unde muchia e(u, v) apartine E daca dist_euclidiana(u, v) <= r
func generateGraph(positions []Coord, n int, r float64) map[int][]int {
	graph := make(map[int][]int)

	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			dist := getEuclideanDistance(positions[i], positions[j])
			if dist <= r {
				graph[i] = append(
					graph[i],
					j,
				)
				graph[j] = append(
					graph[j],
					i,
				)
			}
		}
	}
	return graph
}

func getEuclideanDistance(u Coord, v Coord) float64 {
	xu, xv := u.X, v.X
	yu, yv := u.Y, v.Y

	dx := xv - xu
	dy := yv - yu

	d := math.Sqrt(float64(dx*dx + dy*dy))
	return d
}

func LCM(positions []Coord, node int, in []<-chan Coord, out []chan<- Coord, newPositions []Coord, wg *sync.WaitGroup) {
	defer wg.Done()
	pos := positions[node]

	// trimit pozitia curenta a nodului pe fiecare canal de iesire
	for _, ch := range out {
		ch <- pos
	}

	// new_position = (pos + sum(vecini)) / (1 + len(vecini))
	sumX := pos.X
	sumY := pos.Y

	for _, ch := range in {
		neigh := <-ch
		sumX += neigh.X
		sumY += neigh.Y
	}

	count := float64(len(in) + 1)
	new_pos := Coord{
		sumX / count,
		sumY / count,
	}
	newPositions[node] = new_pos
}

// pot comprima doua perechi de coordonate daca diferenta dintre ele este < 10^(-9)
func equals(u Coord, v Coord) bool {
	epsilon := 1e-9
	return math.Abs(u.X-v.X) < epsilon && math.Abs(u.Y-v.Y) < epsilon
}

func compress(positions []Coord) []Coord {
	newPositions := make([]Coord, 0)

	// verific daca pozitia curenta exista deja (in urma comprimarii) in newPos; altfel o adaug
	for _, posU := range positions {
		found := false
		for _, posV := range newPositions {
			if equals(posU, posV) {
				found = true
				break
			}
		}
		if !found {
			newPositions = append(newPositions, posU)
		}
	}
	return newPositions
}

func main() {
	file, err := os.Open("ex1.txt")
    if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	n, _ := strconv.Atoi(scanner.Text())

	var r float64 = 3.0
	var wg sync.WaitGroup

	positions := make([]Coord, n)
	for i := 0; i < n; i++ {
		scanner.Scan()
		var x, y float64

		// scaneaza linia curenta si mapeaza cele doua valori pe argumentele furnizate, conform formatului
		fmt.Sscanf(scanner.Text(), "%f %f", &x, &y)

		// lista cu coordonatele asociate nodurilor
		positions[i] = Coord{x, y}
	}

	// algoritmul ruleaza pentru maxim 100 de iteratii sau pana cand ramane un singur robot
	for it := 0; it < 100 && n > 1; it++ {
		graph := generateGraph(positions, n, r)
		newPositions := make([]Coord, n)

		in := make([][]<-chan Coord, n)
		out := make([][]chan<- Coord, n)

		// construiesc canalele asociate fiecarei muchii
		for u := 0; u < n; u++ {
			for _, v := range graph[u] {
				ch := make(chan Coord, 1)
				out[u] = append(out[u], ch)
				in[v] = append(in[v], ch)
			}
		}

		for i := 0; i < n; i++ {
			wg.Add(1)
			go LCM(positions, i, in[i], out[i], newPositions, &wg)
		}
		wg.Wait()

		positions = compress(newPositions)
		n = len(positions)
	}

	fmt.Println("Remaining robots: ", n)
	for i, pos := range positions {
		fmt.Printf("Robot %d: (%.9f, %.9f)\n", i, pos.X, pos.Y)
	}
}
