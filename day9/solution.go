package day9

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

type Move struct {
    Dir rune
    Steps int
}

type Pos struct {
    x int
    y int
}

func Part1(input_file string) string {
    moves := parse(input_file)
    knots := [2]Pos {{0,0}}

    position := make(map[int]bool, 0)
    position[knots[1].x << 32 + knots[1].y] = true

    // keep track of unique positions
    for _, move := range moves {
        for i := 0; i < move.Steps; i++ {
            // First move head
            moveHead(&knots[0], move.Dir)
            follow(&knots[1], &knots[0])
            position[knots[1].x << 32 + knots[1].y] = true
        }
    }

    return fmt.Sprint(len(position))
}

func Part2(input_file string) string {
    moves := parse(input_file)
    knots := [10]Pos {{0,0}}

    position := make(map[int]bool, 0)
    position[knots[9].x << 32 + knots[9].y] = true

    // keep track of unique positions
    for _, move := range moves {
        for i := 0; i < move.Steps; i++ {
            // First move head
            moveHead(&knots[0], move.Dir)

            for i := 1; i < len(knots); i++ {
                follow(&knots[i],&knots[i-1])
            }
            position[knots[9].x << 32 + knots[9].y] = true
        }
    }

    return fmt.Sprint(len(position))
}

func moveHead(h *Pos, d rune) {
    switch d {
    case 'R':
        h.x++
    case 'L':
        h.x--
    case 'U':
        h.y--
    case 'D':
        h.y++
    default:
        panic("Invalid Direction")
    }
}

func follow(t *Pos, h *Pos) {
    if abs(h.x - t.x) <= 1 && abs(h.y - t.y) <= 1 { // touching
        return

    } else if abs(h.x - t.x) > 1 && abs(h.y - t.y) == 0 { // Horizontal Difference
        if h.x > t.x {t.x++} else {t.x--}

    } else if abs(h.y - t.y) > 1 && abs(h.x - t.x) == 0 { // Vertical Difference
        if h.y > t.y {t.y++} else {t.y--}

    } else { // One diagonal step
        t.x += norm(h.x - t.x)
        t.y += norm(h.y - t.y)
    }
}

func abs(a int) int {
    if a < 0 {return -a}
    return a
}

func norm(a int) int {
    return a/abs(a)
}


func parse(input_file string) []Move {
    moves := make([]Move, 0)
    f, _ := os.Open(input_file)
    scanner := bufio.NewScanner(f)

    for scanner.Scan() {
        s := scanner.Text()
        dir := rune(s[0])
        n, _ := strconv.Atoi(s[2:])
        moves = append(moves, Move{dir, n})
    }

    return moves
}
