package day18

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Jellycious/aoc2022/utils"
)

type Pos struct {
    x int
    y int
    z int
}

func Part1(input_file string) string {
    positions, grid := parse(input_file)
    sides := 0
    for _, pos := range positions {
        sides += (6 - len(neighbours(pos, &grid)))
    }
    return fmt.Sprint(sides)
}

func Part2(input_file string) string {
    positions, grid := parse(input_file)
    sides := 0

    // Initialize air grid
    air := createAirGrid(&grid)

    for _, pos := range positions {
        ns := neighbours(pos, &air) // gets surrounding air blocks
        sides += len(ns)
    }
    return fmt.Sprint(sides)
}

func createAirGrid(grid *[][][]bool) [][][]bool {
    // Initialize air grid
    air := make([][][]bool, len(*grid))
    for z := range air {
        // Copy plane
        plane := make([][]bool, len((*grid)[0]))
        for y := range plane {
            row := make([]bool, len((*grid)[0][0]))
            plane[y] = row
        }
        air[z] = plane
    }

    // Fill borders with Air
    zlen, ylen, xlen := len(*grid), len((*grid)[0]), len((*grid)[0][0])
    maxlen := utils.Max(utils.Max(zlen, ylen), xlen)
    for i := 0; i < maxlen; i++ {
        for j := 0; j < maxlen; j++ {
            air[utils.Min(zlen-1, i)][utils.Min(ylen-1, j)][0                   ] = true
            air[utils.Min(zlen-1, i)][utils.Min(ylen-1, j)][xlen-1              ] = true
            air[utils.Min(zlen-1, i)][0                   ][utils.Min(xlen-1, j)] = true
            air[utils.Min(zlen-1, i)][ylen-1              ][utils.Min(xlen-1, j)] = true
            air[0                   ][utils.Min(ylen-1, i)][utils.Min(xlen-1, j)] = true
            air[zlen-1              ][utils.Min(ylen-1, i)][utils.Min(xlen-1, j)] = true
        }
    }
    // Fill rest of air
    // 1. Must have a neighbour which is air
    // 2. Must not be lava
    // TODO: FIX THIS UGLY ASS CODE BOYYY
    for z := 1; z < zlen-1; z++ {
        for y := 1; y < ylen-1; y++ {
            for x := 1; x < xlen-1; x++ {
                if (*grid)[z][y][x] { continue } // is lava
                // check whether one of the neighbours is
                ns := neighbours(Pos{x,y,z}, &air)
                if len(ns) > 0 {
                    air[z][y][x] = true
                }
            }
        }
    }

    for z := zlen-2; z > 0; z-- {
        for y := ylen-2; y > 0; y-- {
            for x := xlen-2; x > 0; x-- {
                if (*grid)[z][y][x] { continue } // is lava
                // check whether one of the neighbours is air
                ns := neighbours(Pos{x,y,z}, &air)
                if len(ns) > 0 {
                    air[z][y][x] = true
                }
            }
        }
    }

    for z := 1; z < zlen-1; z++ {
        for y := ylen-2; y > 0; y-- {
            for x := xlen-2; x > 0; x-- {
                if (*grid)[z][y][x] { continue } // is lava
                // check whether one of the neighbours is air
                ns := neighbours(Pos{x,y,z}, &air)
                if len(ns) > 0 {
                    air[z][y][x] = true
                }
            }
        }
    }

    for z := zlen-2; z > 0; z-- {
        for y := ylen-2; y > 0; y-- {
            for x := 1; x < xlen-1; x++ {
                if (*grid)[z][y][x] { continue } // is lava
                // check whether one of the neighbours is air
                ns := neighbours(Pos{x,y,z}, &air)
                if len(ns) > 0 {
                    air[z][y][x] = true
                }
            }
        }
    }

    for z := zlen-2; z > 0; z-- {
        for y := 1; y < ylen-1; y++ {
            for x := xlen-2; x > 0; x-- {
                if (*grid)[z][y][x] { continue } // is lava
                // check whether one of the neighbours is air
                ns := neighbours(Pos{x,y,z}, &air)
                if len(ns) > 0 {
                    air[z][y][x] = true
                }
            }
        }
    }

    for z := zlen-2; z > 0; z-- {
        for y := 1; y < ylen-1; y++ {
            for x := 1; x < xlen-1; x++ {
                if (*grid)[z][y][x] { continue } // is lava
                // check whether one of the neighbours is air
                ns := neighbours(Pos{x,y,z}, &air)
                if len(ns) > 0 {
                    air[z][y][x] = true
                }
            }
        }
    }

    return air
}

func copyGrid(grid *[][][]bool) [][][]bool {
    res := make([][][]bool, len(*grid))
    for z := range res {
        // Copy plane
        plane := make([][]bool, len((*grid)[0]))
        for y := range plane {
            row := make([]bool, len((*grid)[0][0]))
            copy(row[:], (*grid)[z][y])
            plane[y] = row
        }
        res[z] = plane
    }
    return res
}

func printGrid(grid *[][][]bool, inv bool) {
    for z := range *grid {
        fmt.Printf("z = %d:\n", z)
        for y := range (*grid)[0] {
            for x := range(*grid)[0][0] {
                b := (*grid)[z][y][x]
                if inv { b = !b }
                if b {
                    fmt.Print("#")
                }else {
                    fmt.Print(".")
                }
            }
            fmt.Println()
        }
    }
}

func neighbours(pos Pos, grid *[][][]bool) []Pos {
    // Returns neighbours, for which grid[z][y][x] = true
    cands := make([]Pos, 6)
    cands[0] = Pos{pos.x - 1, pos.y, pos.z}
    cands[1] = Pos{pos.x + 1, pos.y, pos.z}
    cands[2] = Pos{pos.x, pos.y - 1, pos.z}
    cands[3] = Pos{pos.x, pos.y + 1, pos.z}
    cands[4] = Pos{pos.x, pos.y, pos.z - 1}
    cands[5] = Pos{pos.x, pos.y, pos.z + 1}

    neighbours := make([]Pos, 0)
    for _, cand := range cands {
        if cand.x >= 0 && cand.y >= 0 && cand.z >= 0 {
            if (*grid)[cand.z][cand.y][cand.x] { neighbours = append(neighbours, cand) }
        }
    }
    return neighbours
}

func parse(input_file string) ([]Pos, [][][]bool) {
    f, _ := os.Open(input_file)
    scanner := bufio.NewScanner(f)
    positions := make([]Pos, 0)
    xmax, ymax, zmax := 0, 0, 0
    for scanner.Scan() {
        l := scanner.Text()
        coords := strings.Split(l, ",")
        x, _ := strconv.Atoi(coords[0])
        xmax = utils.Max(x,xmax)
        y, _ := strconv.Atoi(coords[1])
        ymax = utils.Max(y,ymax)
        z, _ := strconv.Atoi(coords[2])
        zmax = utils.Max(z,zmax)
        positions = append(positions, Pos{x,y,z})

    }

    // Grid is padded with zeroes
    grid := make([][][]bool, zmax+2)
    for z := range grid {
        plane := make([][]bool, ymax+2)
        for y := range plane {
            line := make([]bool, xmax+2)
            plane[y] = line
        }
        grid[z] = plane
    }

    for _, pos := range positions {
        grid[pos.z][pos.y][pos.x] = true
    }

    return positions, grid
}
