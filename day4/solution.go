package day4

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Range struct {
    min int
    max int
}

type RangePair struct {
    left Range
    right Range
}

// Check whether p2 is fully included in p1
func envelops(p1 *Range, p2 *Range) bool {
    return p1.min <= p2.min && p1.max >= p2.max
}

// Checks whether p2 overlaps p1
func isoverlap(p1 *Range, p2 *Range) bool {
    return  (p1.max >= p2.min && p1.min <= p2.max) ||
            (p2.max >= p1.min && p2.min <= p1.max)
}

func Part1(input_file string) string {
    rps := parse(input_file)
    counter := 0
    for _, rp := range rps {
        if envelops(&rp.left, &rp.right) || envelops(&rp.right, &rp.left) {
            counter += 1
        }
    }
    return fmt.Sprint(counter)
}

func Part2(input_file string) string {
    rps := parse(input_file)
    counter := 0
    for _, rp := range rps {
        if isoverlap(&rp.left, &rp.right) {
            counter++
        }

    }
    return fmt.Sprint(counter)
}

// -- PARSING ----------------------------------------------
func parse(input_file string) []RangePair {
    // Read input
    sb, _ := os.ReadFile(input_file)
    s := string(sb)
    pairs := strings.Split(s, "\n") // split on newline
    pairs = pairs[:len(pairs)-1]

    ranges := make([]RangePair, len(pairs)) // allocate array

    for i, p := range pairs {
        pair := strings.Split(p, ",")
        left := strings.Split(pair[0], "-")
        lmin, _ := strconv.Atoi(left[0])
        lmax, _ := strconv.Atoi(left[1])
        right := strings.Split(pair[1], "-")
        rmin, _ := strconv.Atoi(right[0])
        rmax, _ := strconv.Atoi(right[1])
        rp := RangePair{
            Range {lmin, lmax},
            Range {rmin, rmax},
        }
        ranges[i] = rp
    }

    return ranges
}
