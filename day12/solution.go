package day12

import (
	"bufio"
	"fmt"
	"os"

    "github.com/emirpasic/gods/queues/linkedlistqueue"
)


type pos struct {
    x int
    y int
}

type grid [][]byte

func (g grid) String() string {
    // Pretty print the grid
    s := ""
    for y := range g {
        r := g[y]
        for x := range r {
            c := g[y][x]
            s += fmt.Sprintf("%c", c)
        }
        s += fmt.Sprintln()
    }
    return s
}

func findChars(g *grid, c byte) []pos {
    posList := make([]pos, 0)
    for y := range *g {
        for x := range (*g)[0] {
            if (*g)[y][x] == c {
                posList = append(posList, pos{x,y})
            }
        }
    }
    return posList
}

func min(a,b int) int {
    if a > b {return b}
    return a
}

func Part1(input_file string) string {
    g := parse(input_file)
    poss := findChars(&g, 'S')

    bfs := BFS(g, poss[0])

    // Do print magic
    return fmt.Sprint(bfs)
}

func Part2(input_file string) string {
    g := parse(input_file)
    ss := findChars(&g, 'S')
    ss = append(ss, findChars(&g, 'a')...)
    shortest_path := 9999999
    for _,p := range ss {
        nsp := BFS(g, p)
        if nsp > -1 {
            shortest_path = min(shortest_path, BFS(g, p))
        }
    }
    return fmt.Sprint(shortest_path)
}

func BFS(g grid, start pos) int {
    // Finds shortest path through Breadth-first Search.
    // Returns length of the path
    visited := make(map[pos]bool, 0)
    visited[start] = true
    parents := make(map[pos]*pos, 0)

    to_visit := linkedlistqueue.New()
    to_visit.Enqueue(start)

    for {
        cr, ok := to_visit.Dequeue()
        if !ok { // Exhausted search
            return -1
        }

        c := cr.(pos)
        if g[c.y][c.x] == 'E' { // Reached Goal
            // Back trace
            i := 0
            p := &c
            for p != nil {
                p = parents[*p]
                i++
            }
            return i-1
        }

        // Visit neighbours of current node
        for _, n := range neighbours(&g, c) {
            _, present := visited[n]
            if !present {
                // New neighbour
                parents[n] = &c
                to_visit.Enqueue(n)
                visited[n] = true
            }
        }
    }
}


func neighbours(g *grid, p pos) []pos {
    // Return reachable neighbours
    x := p.x
    y := p.y
    c := byteToHeight((*g)[y][x])
    cand := []pos{{x-1,y}, {x+1,y}, {x,y-1}, {x,y+1}}
    neighbours := make([]pos, 0)
    for _, p := range cand {
        n := indexBounded(g, &p)
        if n == -1 {continue} // out of bounds
        if n <= c + 1 { // check whether neighbour is not too high
            neighbours = append(neighbours, p)
        }
    }
    return neighbours
}

func indexBounded(g *grid, p *pos) int {
    if p.y < len(*g) && p.y >= 0 && p.x < len((*g)[0]) && p.x >= 0 {
        return byteToHeight((*g)[p.y][p.x]) // returns item
    }
    return -1 // out of bounds
}

func byteToHeight(b byte) int {
    if b >= 97 { // lower letter
        return int(b-97)
    }else { // capital letter
        if b == 'S' {return 0}
        if b == 'E' {return byteToHeight('z')}
    }
    panic("Shiiiiiit")
}

func parse(input_file string) grid {

    f, _ := os.Open(input_file)
    scanner := bufio.NewScanner(f)

    grid := make([][]byte, 0)

    for scanner.Scan() {

        line := scanner.Text()
        row := make([]byte, len(line))
        for i := range line {
            row[i] = byte(line[i])
        }

        grid = append(grid, row)
    }

    return grid
}
