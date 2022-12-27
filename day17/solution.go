package day17

import (
    "github.com/Jellycious/aoc2022/utils"
	"fmt"
	"os"
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
        time, height = dropRock(&solidified, time, height, rocks[i % len(rocks)], &jet)
    }

    return fmt.Sprint(height+1)
}

func dropRock(solidified *[]Rock, time int, height int, r Rock, jet *Jet) (int, int) {
    /*  Drop a rock until it solidifies
    *   @param solidified: Positions of rocks that have solidified
    *   @param height: Y coordinate of highest piece of rock
    *   @param time: time at start
    *   @param r: Rock to drop
    *   @param jet: the Jetstream that moves rocks around
    *   @returns (newTime, newHeight)
    */

    rock := r.move(Pos{2, height+4}) // offset rocks to start position
    _, rmax := rock.lim()
    curHeight := rmax.y

    done := false

    var newRock, tmpRock Rock
    newRock = rock

    for !done {
        //fmt.Printf("t: %v, h: %v\n", time, curHeight)
        //fmt.Println(newRock)
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
        //fmt.Printf("After jet push %c: %v\n", c, newRock)

        // Move rock down
        tmpRock = newRock // store state
        newRock = newRock.move(Pos{0,-1})
        min, max = newRock.lim()

        if min.y < 0 {
            // Rock has encountered floor
            *solidified = append(*solidified, tmpRock)
            done = true
            curHeight+=1
        } else if collision(solidified, newRock) {
            *solidified = append(*solidified, tmpRock)
            done = true
            curHeight+=1
        }

        //fmt.Printf("After gravity: %v\n", newRock)
        time++
        curHeight--
    }

    return time, utils.Max(height, curHeight)
}

func collision(solidified *[]Rock, r Rock) bool {
    for _, sr := range *solidified {
        if r.collides(sr) { return true }
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
