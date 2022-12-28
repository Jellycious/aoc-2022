package day17

import (
	"fmt"
	"os"
	"reflect"

	"github.com/Jellycious/aoc2022/utils"
)

type Pos struct {
    x int
    y int
}

func(p Pos) add(other Pos) Pos {
    return Pos{p.x + other.x, p.y + other.y}
}

func(p Pos) sub(other Pos) Pos {
    return Pos{p.x - other.x, p.y - other.y}
}

type Rock struct {
    shape []Pos // shape off the rock (a pointer would be better...)
    pos Pos // Position of the rock (every shape is offset by this value)
}

func(r Rock) move(p Pos) Rock {
    return Rock{r.shape, r.pos.add(p)}
}

func(r Rock) lim() (Pos, Pos) {
    // Returns the limits of space the rock touches of the rock
    xmin := int(^uint(0) >> 1)
    xmax := -xmin - 1
    ymin := xmin
    ymax := xmax

    for i := range r.shape {
        xmin = utils.Min(xmin, r.shape[i].add(r.pos).x)
        xmax = utils.Max(xmax, r.shape[i].add(r.pos).x)
        ymin = utils.Min(ymin, r.shape[i].add(r.pos).y)
        ymax = utils.Max(ymax, r.shape[i].add(r.pos).y)
    }
    return Pos{xmin,ymin}, Pos{xmax,ymax}
}

func(r Rock) collides(other Rock) bool {
    // Checks whether two rocks have a collision
    for i := range r.shape {
        p := r.shape[i].add(r.pos)
        for j := range other.shape {
            p2 := other.shape[j].add(other.pos)
            if p == p2 {return true}
        }
    }
    return false
}

type Jet []byte

func Part1(input_file string) string {
    rocks := getRocklist()
    jet := parse(input_file)
    solidified := make([]Rock, 0)

    time, height := 0, -1
    for i := 0; i < 2022; i++ {
        time, height, _ = dropRock(&solidified, time, height, rocks[i % len(rocks)], &jet)
    }

    return fmt.Sprint(height+1)
}

// State
type State struct {
    jetI int
    rockI int
    /* Offset of spawn position of rock and final position
       Is somewhat of an approximate state, so might not work for every input */
    prevOffsets []Pos
    /* Misc Data */
    i int // rock number
    height int // height after rock had been dropped
}

func(s State) peq(other State) bool {
    // Partial equality for states (excludes Misc Data in State)
    if s.jetI != other.jetI || s.rockI != other.rockI { return false }
    if reflect.DeepEqual(s.prevOffsets, other.prevOffsets) {
        return true
    }
    return false
}

func Part2(input_file string) string {
    /* Optimization ideas: Cycle Detection*/
    rocks := getRocklist()
    jet := parse(input_file)
    solidified := make([]Rock, 0)
    prevStates := make([]State, 0)

    PREV_SIZE := 30 /* Increase this value for more accurate results */
    prevOffsets := make([]Pos, PREV_SIZE + 1) // Keep track of previous 500 offsets

    i, time, height := 0, 0, -1

    var cycleStart, cycleRepeat State

    for {
        oldTime := time

        var offset Pos
        time, height, offset = dropRock(&solidified, time, height, rocks[i % len(rocks)], &jet)

        // Copy new offset into prevOffsets and shift list
        prevOffsets[PREV_SIZE] = offset
        offsets := make([]Pos, PREV_SIZE)
        copy(prevOffsets[:PREV_SIZE], prevOffsets[1:PREV_SIZE+1])
        copy(offsets, prevOffsets[:PREV_SIZE])

        newState := State {
            oldTime % len(jet),
            i % len(rocks),
            offsets,
            i,
            height,
        }

        // Cycle detection
        s := contains(&prevStates, newState)
        if s != nil {
            cycleStart = *s
            cycleRepeat = newState
            i++
            break
        }

        // add new state to previous states
        prevStates = append(prevStates, newState)
        i++
    }

    // Compute the height after 1000000000000 rocks using the detected cycle
    cycleRocks := cycleRepeat.i - cycleStart.i // Rocks that get dropped in a single cycle
    cycleHeight := cycleRepeat.height - cycleStart.height // The height that gets added in a single cycle

    rocksRemaining := 1000000000000 - cycleRepeat.i // Rocks that should still be dropped after the cycle has been detected
    cyclesRemaining := rocksRemaining / cycleRocks // how many cycles can we fit in the remaining rocks
    lastRocks := rocksRemaining % cycleRocks // last few computations that are required

    additionalHeight := cycleHeight * cyclesRemaining // Height accumulated in the remaining cycles

    // Simulate last few rocks
    for j := 0; j < lastRocks - 1; j++ {
        time, height, _ = dropRock(&solidified, time, height, rocks[i % len(rocks)], &jet)
        i++
    }

    return fmt.Sprint(height+additionalHeight+1)
}

func contains(states *[]State, state State) *State {
    for _, s := range *states {
        if state.peq(s) { return &s }
    }
    return nil
}

func dropRock(solidified *[]Rock, time int, height int, r Rock, jet *Jet) (int, int, Pos) {
    /*  Drop a rock until it solidifies
    *   @param solidified: Positions of rocks that have solidified
    *   @param height: Y coordinate of highest piece of rock
    *   @param time: time at start
    *   @param r: Rock to drop
    *   @param jet: the Jetstream that moves rocks around
    *   @returns (newTime, newHeight, offset)
    */

    rock := r.move(Pos{2, height+4}) // offset rocks to start position
    _, rmax := rock.lim()
    curHeight := rmax.y

    done := false

    var newRock, tmpRock Rock
    newRock = rock

    var endPos, startPos Pos
    startPos = rock.pos

    for !done {
        tmpRock = newRock
        c := (*jet)[time % len(*jet)]

        // Push rock left or right
        if c == '<' { newRock = newRock.move(Pos{-1,0}) } else { newRock = newRock.move(Pos{1,0}) }

        // verify that push did not result in a collision
        min, max := newRock.lim()
        if min.x < 0 || max.x > 6 {
            newRock = tmpRock // revert to old non pushed state
        } else if collision(solidified, newRock) {
            newRock = tmpRock
        }

        // Move rock down
        tmpRock = newRock // store state
        newRock = newRock.move(Pos{0,-1})
        min, max = newRock.lim()

        if min.y < 0 {
            // Rock has encountered floor
            *solidified = append(*solidified, tmpRock)
            endPos = tmpRock.pos
            done = true
            curHeight+=1
        } else if collision(solidified, newRock) {
            *solidified = append(*solidified, tmpRock)
            endPos = tmpRock.pos
            done = true
            curHeight+=1
        }

        time++
        curHeight--
    }

    return time, utils.Max(height, curHeight), endPos.sub(startPos)
}

func collision(solidified *[]Rock, r Rock) bool {
    BOUND := 20 // Bound (**Increase this value for more accurate results */
    for i := len(*solidified) - 1; i >= utils.Max(0, (len(*solidified) - BOUND)); i-- {
        if (*solidified)[i].collides(r) { return true }

    }

    return false
}

func parse(input_file string) Jet {
    bs, _ := os.ReadFile(input_file)
    jet := Jet(make([]byte, len(bs)-1))

    for i := range jet {
        jet[i] = bs[i]
    }

    return jet
}

func getRocklist() []Rock {
    rocks := make([]Rock, 5)
    rocks[0] = Rock{[]Pos{{0,0},{1,0},{2,0},{3,0}}, Pos{0,0}}
    rocks[1] = Rock{[]Pos{{1,2},{0,1},{1,1},{2,1},{1,0}}, Pos{0,0}}
    rocks[2] = Rock{[]Pos{{2,2},{2,1},{0,0},{1,0},{2,0}}, Pos{0,0}}
    rocks[3] = Rock{[]Pos{{0,3},{0,2},{0,1},{0,0}}, Pos{0,0}}
    rocks[4] = Rock{[]Pos{{0,1},{1,1},{0,0},{1,0}}, Pos{0,0}}
    return rocks
}
