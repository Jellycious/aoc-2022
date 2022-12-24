package day16

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/emirpasic/gods/queues/linkedlistqueue"
)

type Valve struct {
    id string
    rate int
    neighbours []string
}

type Flow struct {
    time_open int
    rate int
}

type FlowMap map[string]Flow

func (m FlowMap) totalFlow() int {
    total := 0
    for _, v := range m {
        total += v.rate * v.time_open
    }
    return total
}

func Part1(inputFile string) string {
    graph := parse(inputFile)

    closed := make([]string, 0)
    for id := range graph {
        if graph[id].rate > 0 {
            closed = append(closed, id)
        }
    }

    travelMap := createTravelMap(graph)
    res := DFS("AA", 30, closed, &travelMap, &graph)
    return fmt.Sprint(res.totalFlow())
}


func createTravelMap(graph map[string]Valve) map[string]map[string]int {
    // Returns the lenghts of the shortest paths from nodes
    pois := make([]string, 0)
    for id := range graph {
        if graph[id].rate > 0 {
            pois = append(pois, id)
        }
    }

    travelMap := make(map[string]map[string]int, 0)
    travelMap["AA"] = BFS("AA", pois, graph)
    for _, poi := range pois {
        travelMap[poi] = BFS(poi, pois, graph)
        delete(travelMap[poi], poi)
    }

    return travelMap
}

type Node struct {
    v Valve
    p *Node
}

func BFS(valveId string, ids []string, graph map[string]Valve) map[string]int {
    // returns length of shortest paths from valveId to valves in ids
    startValve := graph[valveId]
    start := Node{startValve, nil}
    visited := make(map[string]struct{})
    visited[startValve.id] = struct{}{}

    eta := make(map[string]int, len(ids))

    queue := linkedlistqueue.New()
    queue.Enqueue(start)

    for !queue.Empty() {
        nq, ok := queue.Dequeue()
        if !ok {break} // Exhausted graph

        n := nq.(Node)

        if contains(ids, n.v.id) {
            // Found a PoI, store length
            parent := n.p
            length := 0

            for parent != nil {
                length++
                parent = parent.p
            }

            eta[n.v.id] = length + 1 // +1 because we always open valves (time incoorperated in path length)
        }

        // Visit neighbours
        for _, nId := range n.v.neighbours {
            _, visitedBefore := visited[nId]
            if !visitedBefore {
                visited[nId] = struct{}{}
                queue.Enqueue(Node{graph[nId], &n})
            }
        }
    }

    for _, id := range ids {
        _, check := eta[id]
        if !check {
            fmt.Printf("Could not reach %v from %v\n", id, valveId)
            panic("Acyclic Graph!")

        }
    }

    return eta
}


func DFS(valveId string, time int, closed []string, travelMap *map[string]map[string]int, graph *map[string]Valve) FlowMap {

    // Open this valve
    valve := (*graph)[valveId]

    flowMap := FlowMap(make(map[string]Flow, 0))

    for i, id := range closed {
        // Check whether we can reach and if so visit that node
        pathLength := (*travelMap)[valveId][id]

        if pathLength < time { // we can reach this node and open it within timespan
            new_closed := make([]string, len(closed)-1)
            copy(new_closed, closed[:i])
            copy(new_closed[i:], closed[i+1:])

            bound := computeUpperBound(time - pathLength, new_closed, graph) + (time - pathLength) * (*graph)[id].rate

            if flowMap.totalFlow() <= bound {
                candFlow := DFS(id, time - pathLength, new_closed, travelMap, graph) // we might want to do boundcheck before going into DFS

                if candFlow.totalFlow() > flowMap.totalFlow() { // New better optimal solution.
                    flowMap = candFlow
                }
            }
        }
    }

    flowMap[valveId] = Flow{time, valve.rate}
    return flowMap
}

func computeUpperBound(time int, closed []string, graph *map[string]Valve) int {
    valve_rates := make([]int, len(closed))
    for k := range closed {
        valve_rates[k] = (*graph)[closed[k]].rate
    }
    sorted := sort.IntSlice(valve_rates)
    sort.Sort(sorted)
    total := 0
    for _, v := range sorted {
        if time < 0 { break }
        time = time - 2 // time to reach, time to open
        total += v * time
    }
    return total
}

func contains[T comparable](l []T, e T) bool {
    for _, elem := range l {
        if elem == e {
            return true
        }
    }
    return false
}


// -- PARSING ----------------------------------------------
func parse(inputFile string) map[string]Valve {
    f, _ := os.Open(inputFile)
    scanner := bufio.NewScanner(f)
    res := make(map[string]Valve, 0)

    for scanner.Scan() {
        l := scanner.Text()
        valveId := l[6:8]
        ns := strings.Split(l, "=")
        ns = strings.Split(ns[1], ";")
        rate, _ := strconv.Atoi(ns[0])
        neighbours := strings.Split(ns[1], ", ")
        prefixL := len(neighbours[0])
        neighbours[0] = neighbours[0][prefixL-2:prefixL]

        v := Valve {
            valveId,
            rate,
            neighbours,
        }
        res[valveId] = v
    }
    return res
}

