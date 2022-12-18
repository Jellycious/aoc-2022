package day15

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

type Range struct {
    s int
    e int // exclusive
}


func (r Range) nogap(other Range) bool {
    // Returns true if there is no empty space between union of the two ranges.
    nmin := max(r.s, other.s)
    nmax := min(r.e, other.e)
    if nmin > nmax {
        return false
    }
    return true
}

func (r Range) union(other Range) Range {
    if !r.nogap(other) { panic("Cannot unite two ranges that do have a gap.") }
    return Range{min(r.s, other.s), max(r.e, other.e)}
}

func (r Range) less(other Range) bool {
    return r.s < other.s
}

type Sensor struct {
    pos Pos
    beacon Pos
    perim int
}

func Part1(input_file string) string {
    sensors := parse(input_file)

    y := 2000000

    ranges := make([]Range, 0)

    rsmin := int(^uint(0)>>1)
    remax := - rsmin - 1

    beacons_on_row := make(map[Pos]struct{}, 0)

    for _, s := range sensors {
        // Check whether y lies in its perimeter
        if s.beacon.y == y {
            beacons_on_row[s.beacon] = struct{}{}
        }
        if y >= s.pos.y - s.perim && y <= s.pos.y + s.perim {
            // include range at this y
            dy := abs(s.pos.y - y)
            rs := s.pos.x - (s.perim - dy)
            rsmin = min(rsmin, rs)
            re := s.pos.x + (s.perim - dy) + 1 // e not inclusive thus +1
            ranges = append(ranges, Range{rs,re})
            remax = max(remax, re)
        }
    }

    row := make([]bool, remax - rsmin)
    for _, r := range ranges {
        for i := r.s; i < r.e; i++ {
            row[i - rsmin] = true
        }
    }

    count := 0
    for i := range row {
        if row[i] { count++ }
    }

    for range beacons_on_row {
        count--
    }

    return fmt.Sprint(count)
}

func Part2(input_file string) string {
    sensors := parse(input_file)

    LIM := 4000000
    //LIM := 20

    ranges := make([][]Range, LIM+1)
    // initialize empty array for every row
    for i := range ranges {
        ranges[i] = make([]Range, 0)
    }

    for _, s := range sensors {
        // Check whether y lies in its perimeter

        for y := max(0, s.pos.y - s.perim); y < min(LIM, s.pos.y + s.perim) + 1; y++ {
            // include range at this y
            dy := abs(s.pos.y - y)
            rs := s.pos.x - (s.perim - dy)
            re := s.pos.x + (s.perim - dy) + 1 // e not inclusive thus +1
            r := Range{rs,re}
            ranges[y] = insert(r, ranges[y])

        }
    }

    // Go through every row and check whether a gap exists
    answer := 0
    for i := 0; i < LIM+1; i++ {
        row := ranges[i]

        n := row[0]
        for j := 1; j < len(row); j++ {
            if n.nogap(row[j]) {
                n = n.union(row[j])
            }else {
                answer = min(n.e, row[j].e) * 4000000 + i
                break
            }
        }
    }
    return fmt.Sprint(answer)
}

func insert(r Range, rs []Range) []Range {
    // See where we can insert something
    for i := range rs {
        ro := rs[i]
        if r.nogap(ro) {
            rs[i] = r.union(ro)
            return rs;
        }
        if r.less(ro) {
            // Insert r at rs[i]
            rs = append(rs, r)
            copy(rs[i+1:], rs[i:])
            rs[i] = r
            return rs;
        }
    }
    // Exhausted
    rs = append(rs, r)
    return rs
}

func abs(a int) int {
    if a < 0 {return -a}
    return a
}

func dist(p1, p2 Pos) int {
    return abs(p1.x - p2.x) + abs(p1.y - p2.y)
}

func min(a, b int) int {
    if a < b { return a }
    return b
}

func max(a, b int) int {
    if a > b { return a }
    return b
}

func parse(input_file string) []Sensor {
    f, _ := os.Open(input_file)
    scanner := bufio.NewScanner(f)

    sensors := make([]Sensor, 0)
    for scanner.Scan() {
        l := scanner.Text()
        ss := strings.Split(l, ": closest beacon is at ")
        s := strings.Split(ss[0][10:], ", ")
        b := strings.Split(ss[1], ", ")
        sx, _ := strconv.Atoi(s[0][2:])
        sy, _ := strconv.Atoi(s[1][2:])
        bx, _ := strconv.Atoi(b[0][2:])
        by, _ := strconv.Atoi(b[1][2:])
        ps := Pos{sx,sy}
        pb := Pos{bx,by}
        perim := dist(ps, pb)

        sensors = append(sensors, Sensor{ps, pb, perim})
    }

    return sensors
}
