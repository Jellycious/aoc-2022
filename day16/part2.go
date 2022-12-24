package day16

import (
	"fmt"
	"sort"

	"golang.org/x/exp/constraints"
)


func Part2(inputFile string) string {
    // Lord save me 
    graph := parse(inputFile)

    pois := make([]string, 0)
    pois = append(pois, "AA") // append "AA" to pois
    mapping := make(map[string]int, 0) // maps string to index in pois
    mapping["AA"] = 0
    for id := range graph {
        if graph[id].rate > 0 {
            pois = append(pois, id)
            mapping[id] = len(pois) - 1 // enumeration of ids in pois, pois[index-1] = id
        }
    }

    travelMap := createTravelMap(graph)


    closed := (^(uint32(0)) << len(pois)) | 0x01 // Set closed valves to zero
    available := (^(uint32(0)) << len(pois)) | 0x01 // Set closed valves to zero
    startState := State {
        26,
        closed,
        available,
        0,
        0,
    }
    memoization := make(map[State]int, 0)
    result := BranchAndBound(startState, &travelMap, &graph, &pois, &memoization)

    return fmt.Sprint(result)
}


type State struct {
    time uint32 // time left to do things
    closed uint32 //bit mask of closed valves '1' = open, '0' = closed
    available  uint32 // bit mask of valves which are available '1' = unavailable, '0' = available (A valve becomes unavailable when someone is heading towards it)
    // bitmasks for entity me and elephant
    // bitmask contains current location, eta and next destination b23..b16 -> location, b15..8 -> time, b7..b0 -> dest
    // if an entity is not allowed to make a choice it will be set to 0xffffff = ^uint32(0)
    me uint32
    elephant uint32
}

func (s State) String(pois []string) string {
    r := "State:\n"
    r += "Time Left: " + fmt.Sprint(s.time) + "\n"
    r += "Open:        " + valveBitmaskPrettyPrint(s.closed, pois) + "\n"
    r += "Unavailable: " + valveBitmaskPrettyPrint(s.available, pois) + "\n"
    r += "Me: Position=" + fmt.Sprint(bmGetPos(s.me)) + ", Dst=" + fmt.Sprint(bmGetDst(s.me)) + ", Eta=" + fmt.Sprint(bmGetEta(s.me)) + "\n"
    r += "Elephant: Position=" + fmt.Sprint(bmGetPos(s.elephant)) + " Dst=" + fmt.Sprint(bmGetDst(s.elephant)) + ", Eta=" + fmt.Sprint(bmGetEta(s.elephant)) + "\n"
    return r
}

func (s State) totalFlowRate(graph *map[string]Valve, pois *[]string) uint {

    flow := uint(0)
    for i := uint32(0); i < uint32(len(*pois)); i++ {
        if !bmCheckBit(s.closed, i) {
            flow += uint((*graph)[(*pois)[i]].rate)
        }
    }
    return flow
}

// Bitmask functions
func bmGetEta(bitmask uint32) uint32 {
    return (bitmask >> 8) & 0xff
}

func bmGetDst(bitmask uint32) uint32 {
    return bitmask & 0xff
}

func bmGetPos(bitmask uint32) uint32 {
    return (bitmask >> 16) & 0xff
}

func bmSetEta(bitmask uint32, eta uint32) uint32 {
    return (bitmask & 0xff00ff)  | (eta << 8)
}

func bmSetDst(bitmask uint32, dst uint32) uint32 {
    return (bitmask & 0xffff00) | dst
}

func bmSetPos(bitmask uint32, pos uint32) uint32 {
    return (bitmask & 0x00ffff) | (pos << 16)
}

func bmCheckBit(bitmask uint32, index uint32) bool {
    return (bitmask >> index) & 0x01 == 0
}

func bmCloseValve(bitmask uint32, index uint32) uint32 {
    return bitmask | (0x01 << index)
}

func valveBitmaskPrettyPrint(closed uint32, pois []string) string {
    s := ""
    for i := 0; i < len(pois); i++ {
        if (closed >> i) & 0x01 == 0 {
            s += pois[i] + ":[ ] "
        }else {
            s += pois[i] + ":[X] "
        }
    }
    return s
}

// Gets all next possible states from the current state 
// It already does prune the search space significantly, hence the size of this function
func (s State) neighbours(travelMap *map[string]map[string]int, pois *[]string) []State {
    time := s.time
    closed := s.closed

    me := s.me

    elephant := s.elephant

    if bmGetEta(me) != 0 && bmGetEta(elephant) != 0 {
        return make([]State, 0)
    }

    if time <= 0 || closed == (^uint32(0)) { return make([]State, 0) } // No new states possible (time exhausted or all valves are open

    if bmGetEta(me) == 0 && bmGetEta(elephant) != 0 {
        nextStates := s.moveEntity(false, travelMap, pois)
        for i := range nextStates {
            nextStates[i].timeWalk()
        }
        return nextStates
    }

    if bmGetEta(elephant) == 0 && bmGetEta(me) != 0 {
        nextStates := s.moveEntity(true, travelMap, pois)
        for i := range nextStates {
            nextStates[i].timeWalk()
        }
        return nextStates
    }

    var newStates []State

    // Both elephant and me can make a move
    if bmGetEta(me) == 0 && bmGetEta(elephant) == 0 {
        elephantMoves := s.moveEntity(true, travelMap, pois)
        meMoves := s.moveEntity(false, travelMap, pois)
        combinedChoices := make([]State, 0)

        for _, s := range meMoves {
            if s.available == ^uint32(0) {
                combinedChoices = append(combinedChoices, s)
                continue
            }

            combinedChoices = append(combinedChoices, s.moveEntity(true, travelMap, pois)...)
        }
        for _, s := range elephantMoves {
            if s.available == ^uint32(0) {
                combinedChoices = append(combinedChoices, s)
                continue
            }

            combinedChoices = append(combinedChoices, s.moveEntity(false, travelMap, pois)...)
        }
        newStates = combinedChoices
    }

    for i := range newStates {
        newStates[i].timeWalk()
    }

    // There will be duplicate states
    newStates = removeDuplicate(newStates)

    return newStates
}

// Makes all possible choices for an entity possible from state
// It returns new state of the same time, but with entity and valves set accordingly to the corresponding choice.
// For example, if me goes to valve "BB", then it will set the destination and eta time for me. Furthermore, it will set the valve "BB" to unavailable.
func (state *State) moveEntity(elephant bool, travelMap *map[string]map[string]int, pois *[]string) []State {
    entity := state.me
    nextStates := make([]State, 0)

    if elephant { entity = state.elephant }

    // Entity is not allowed to make any moves
    if entity  == ^uint32(0) || bmGetEta(entity) != 0 {
        fmt.Printf("Entity is elephant? %v\n", elephant)
        fmt.Println(state.String(*pois))
        panic("Can't move!!!")

    } // return empty list

    available := state.available
    time := state.time
    pos := bmGetPos(entity)

    if available == ^uint32(0) {
        if elephant {
            return []State{{state.time, state.closed, state.available, state.me, ^uint32(0)}}
        }
        return []State{{state.time, state.closed, state.available, ^uint32(0), state.elephant}}
    } // all valves are unavailable so no new move is made.

    for i := uint32(0); i < uint32(len(*pois)); i++ {

        if bmCheckBit(available, i) { // valve at index 'i' in pois is available
            // Me will travel to this new valve
            time_to_dst := uint32((*travelMap)[(*pois)[pos]][(*pois)[i]])
            if time_to_dst <= time { // we can reach this location in time
                // Set new destination and time to destination for entity
                newEntity := bmSetEta(entity, time_to_dst)
                newEntity = bmSetDst(newEntity, uint32(i))
                // Set the valve that entity is traveling to as unavailable
                newAvailable := available | (0x01 << uint32(i))

                if elephant {
                    nextStates = append(nextStates, State {state.time, state.closed, newAvailable, state.me, newEntity})
                }else {
                    nextStates = append(nextStates, State {state.time, state.closed, newAvailable, newEntity, state.elephant})
                }
            }
        }
    }

    if len(nextStates) == 0 {
        // No new state possible, so return original state, but with entity disabled (entity = ^uint32(0))
        nextStates = []State{{state.time, state.closed, available, ^uint32(0), state.elephant}}
        if elephant { nextStates = []State{{state.time, state.closed, available, state.me, ^uint32(0)}} }
    }

    return nextStates
}

// Moves time forward until a new choice needs to be made.
func (s *State) timeWalk() {
    meEta := bmGetEta(s.me)
    elephantEta := bmGetEta(s.elephant)

    dt := min(meEta, elephantEta)

    if dt == 255 || dt > s.time { // Time cannot be passed forward
        return
    }

    newTime := s.time - dt
    newClosed := s.closed

    newMeEta := meEta
    if newMeEta < 255 { newMeEta = newMeEta - dt }

    newMe := bmSetEta(s.me, newMeEta) // update eta for me

    if newMeEta == 0 {
        newPos := bmGetDst(s.me)
        // Update me 
        newMe = bmSetDst(newMe, 0)
        newMe = bmSetPos(newMe, newPos)
        newClosed = bmCloseValve(newClosed, newPos)
    }

    newElephantEta := elephantEta
    if newElephantEta < 255 { newElephantEta = newElephantEta - dt }

    newElephant := bmSetEta(s.elephant, newElephantEta) // update eta for me

    if newElephantEta == 0 {
        newPos := bmGetDst(s.elephant)
        // Update elephant 
        newElephant = bmSetDst(newElephant, 0)
        newElephant = bmSetPos(newElephant, newPos)
        newClosed = bmCloseValve(newClosed, newPos)
    }

    s.time = newTime
    s.closed = newClosed
    s.me = newMe
    s.elephant = newElephant
}

type BBNode struct {
    s State
    bound int
}

type NodeList []BBNode

func (nl NodeList) Less(i, j int) bool {
    return nl[i].bound < nl[j].bound
}

func (nl NodeList) Len() int { return len(nl) }

func (nl NodeList) Swap(i, j int) {
    tmp := nl[i]
    nl[i] = nl[j]
    nl[j] = tmp
}


// Returns maximum flow for this branch
// You know exactly what flow is counted double because you have this information
func BranchAndBound(state State, travelMap *map[string]map[string]int, graph *map[string]Valve, pois *[]string, memoization *map[State]int) int {

    currentFlowRate := state.totalFlowRate(graph, pois)
    currentFlow := int(currentFlowRate) * int(state.time)

    // Get all new states
    nextStates := state.neighbours(travelMap, pois)
    nodes := NodeList(make([]BBNode, len(nextStates)))

    for i, s := range nextStates {
        b := computeUpperBound2(&s, pois, graph)
        nodes[i] = BBNode{s, b}
    }

    sort.Sort(sort.Reverse(nodes))

    var bestFlow int = currentFlow

    for _, node := range nodes {

        if node.bound < bestFlow { continue } // Upper bound is smaller than current best

        // Memoization
        tflow, ok := (*memoization)[node.s]
        if !ok {
            tflow = BranchAndBound(node.s, travelMap, graph, pois, memoization)
            (*memoization)[node.s] = tflow
        }

        // Check if this flow is an improvement
        actualFlow := currentFlow + tflow - (int(node.s.time) * int(currentFlowRate))
        bestFlow = max(bestFlow, actualFlow)
    }

    return bestFlow
}

func computeUpperBound2(state *State, pois *[]string, graph *map[string]Valve) int {
    // shortestDist is map of the shortest possible distance to reach a valve (from any node)
    // Returns upper bound of flow that this state can generate from state.time
    // Can be improved...
    total := int(state.totalFlowRate(graph, pois)) * int(state.time)
    // check current flow and include this in measurement

    valveRates := make([]int, 16) // Solution is hard-coded to work only with 16 valves (including "AA")

    for k := range *pois {
        v := (*graph)[(*pois)[k]]
        if bmCheckBit(state.available, uint32(k)) {
            valveRates[k] = v.rate
        }

        // If valve is open we know the rate
        if !bmCheckBit(state.closed, uint32(k)) {
           total += int(v.rate) * int(state.time)
        }
    }


    sorted := sort.IntSlice(valveRates)
    sort.Sort(sorted)

    time := state.time
    meEta := bmGetEta(state.me)
    elephantEta := bmGetEta(state.elephant)

    if meEta == 0 { meEta = 2 }
    if elephantEta == 0 { elephantEta = 2 }

    i := 0
    for time > 0 && i < len(valveRates) {
        if meEta == 0 {
            total += valveRates[i] * int(time)
            i++
            meEta = 2
        }

        if i >= len(valveRates) { break }

        if elephantEta == 0 {
            total += valveRates[i] * int(time)
            i++
            elephantEta = 2
        }

        // Time happens
        meEta--
        elephantEta--
        time--
    }

    return total
}


func removeDuplicate(sliceList []State) []State {
    allKeys := make(map[State]bool)
    list := []State{}
    for _, item := range sliceList {
        if _, value := allKeys[item]; !value {
            allKeys[item] = true
            list = append(list, item)
        }
    }
    return list
}

func min[T constraints.Ordered](a, b T) T {
    if a < b { return a }
    return b
}

func max[T constraints.Ordered](a, b T) T {
    if a > b { return a }
    return b
}
