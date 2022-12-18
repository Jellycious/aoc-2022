package day14

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Pos struct {
    x int
    y int
}

const (
    Air int = iota
    Rock
    Sand
    Void
)

type Grid [][]int

func (g Grid) String() string {
    s := ""
    for r := range g {
        for c := range g[r] {
            switch g[r][c] {
            case Air:
                s += "."
            case Rock:
                s += "#"
            case Sand:
                s += "o"
            case Void:
                s += "~"
            }
        }
        s += "\r\n"
    }
    return s
}

func Part1(input_file string) string {
    pr := parse(input_file)
    xmin := pr.xmin
    grid := pr.g
    path := make([]Pos, 1)
    path[0] = Pos{500, 0}

    answer := 0
    for {
        stable := extendPath(&grid, &path, xmin)
        if stable { break }
        answer++ // we add a new sand element

        // Last element becomes Sand
        pos := path[len(path)-1]
        grid[pos.y][pos.x-xmin] = Sand
        path = path[:len(path)-1]
    }

    // Modify grid to have path visible
    for i := range path {
        pos := path[i]
        grid[pos.y][pos.x-xmin] = Void
    }

    return fmt.Sprint(answer)
}

func Part2(input_file string) string {
    pr := parse(input_file)
    grid := pr.g
    path := make([]Pos, 1)
    path[0] = Pos{500, 0}

    // Add two rows to grid
    grid = append(grid, make([]int, len(grid[0])))
    grid = append(grid, make([]int, len(grid[0])))
    // Extend grid on the right and left
    xmin := pr.xmin - (len(grid) - 1)
    for ri := range grid {
        left := make([]int, len(grid)-1)
        right := make([]int, len(grid)-1)
        grid[ri] = append(left, grid[ri]...)
        grid[ri] = append(grid[ri], right...)
    }
    // Fill last row with blocks
    for i := range grid[0] {
        grid[len(grid)-1][i] = Rock
    }

    answer := 0
    for {
        if len(path) == 0 {
            break
        }

        extendPath(&grid, &path, xmin)
        answer++ // we add a new sand element

        // Last element becomes Sand
        pos := path[len(path)-1]
        grid[pos.y][pos.x-xmin] = Sand
        path = path[:len(path)-1]

    }

    // Modify grid to have path visible
    for i := range path {
        pos := path[i]
        grid[pos.y][pos.x-xmin] = Void
    }

    return fmt.Sprint(answer)
}

func extendPath(grid *Grid, p *[]Pos, xmin int) bool {
    // Extends path until it reaches a position where Sand comes to rest
    // Returns true when path ends in Void
    g := *grid
    pos := (*p)[len(*p)-1]

    for {
        // Check for void underneath
        if pos.y + 1 >= len(g) {
            return true
        }
        // Check underneath
        px := pos.x - xmin //offset for x
        if g[pos.y+1][px] == Air {
            // We can go down
            new_pos := Pos{px+xmin, pos.y+1}
            *p = append(*p, new_pos)
            pos = new_pos
            continue
        }

        //Check left down
        if px - 1 < 0 {
            return true
        }

        if g[pos.y+1][px-1] == Air {
            new_pos := Pos{px-1+xmin, pos.y+1}
            *p = append(*p, new_pos)
            pos = new_pos
            continue
        }

        // Check right down
        if px + 1 >= len(g[0]) {
            return true
        }

        if g[pos.y+1][px+1] == Air {
            new_pos := Pos{px+1+xmin, pos.y+1}
            *p = append(*p, new_pos)
            pos = new_pos
            continue
        }

        // Exhausted our options
        break
    }

    return false
}

type parseRes struct {
    g Grid
    xmin int // We offset by this value to get a range of {0, xmax-xmin}
}

func parse(input_file string) parseRes {
    f, _ := os.Open(input_file)
    scanner := bufio.NewScanner(f)

    xmin := int(^uint(0) >> 1)
    xmax := 0
    ymax := 0

    paths := make([][]Pos, 0)
    for scanner.Scan() {
        l := scanner.Text()
        coords := strings.Split(l, " -> ")

        poss := make([]Pos, 0)
        for i := range coords {
            c := parseCoord(coords[i])
            poss = append(poss, c)
            // Keep track coordinate offsets
            xmin = min(xmin, c.x)
            xmax = max(xmax, c.x)
            ymax = max(ymax, c.y)
        }
        paths = append(paths, poss)
    }

    // Construct the grid from the paths
    // Initialize Grid
    g := Grid(make([][]int, 0))
    for y := 0; y < ymax + 1; y++ {
        r := make([]int, (xmax - xmin) + 1)
        g = append(g, r)
    }

    // Fill Grid
    for pi := range paths {
        p := paths[pi]
        for i := 1; i < len(p); i++ {
            fillLine(&g, &p[i-1], &p[i], xmin, Rock)
        }
    }

    return parseRes{Grid(g), xmin}
}

func fillLine(g *Grid, p1 *Pos, p2*Pos, xmin, e int) {
    // Fills the grid between p1 and p2 with e

    dx := norm(p2.x - p1.x)
    dy := norm(p2.y - p1.y)
    // Fill horizontally
    x := p1.x
    y := p1.y
    for x != (p2.x+dx) || y != (p2.y+dy) {
        (*g)[y][x - xmin] = e
        x = x + dx
        y = y + dy
    }
}

func parseCoord(s string) Pos {
    ss := strings.Split(s, ",")
    x, e1 := strconv.Atoi(ss[0])
    y, e2 := strconv.Atoi(ss[1])
    if e1 != nil || e2 != nil { fmt.Println(ss); panic("Conversion Error!") }
    return Pos{x,y}
}

func min(a, b int) int {
    if a < b {return a}
    return b
}

func max(a, b int) int {
    if a > b {return a}
    return b
}

func abs(a int) int {
    if a < 0 { return -a }
    return a
}

func norm(a int) int {
    if a < 0 { return -1 }
    if a > 0 { return 1 }
    return 0
}

