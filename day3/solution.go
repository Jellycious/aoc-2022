package day3

import (
	"bufio"
	"fmt"
	"os"
)

func priority(b byte) int {
    if b >= 97 {return int(b) - 96}else {return int(b) - 38}
}

func contains(l []int, v int) bool {
    for _, vl := range l {
        if vl == v { return true }
    }
    return false
}

func Part1(input_file string) string {
    ruckSacks := parse(input_file)

    priority_sum := 0
    for _, ruckSack := range ruckSacks {
        // convert elements to priorities
        a := make([]int, len(ruckSack))
        for i, v := range ruckSack {
            a[i] = priority(v)
        }

        left := a[:len(a)/2]
        right := a[len(a)/2:]
        // find duplicates
        for _, vl := range left {
            if contains(right, vl) {
                priority_sum += vl
                break;
            }
        }
    }

    return fmt.Sprint(priority_sum)
}

func Part2(input_file string) string {
    ruckSacks := parse(input_file)
    priority_sum := 0


    // convert to chars to priorities
    priorities := make([][]int, len(ruckSacks))
    for i := 0; i < len(ruckSacks); i++ {
        priorities[i] = make([]int, len(ruckSacks[i]))
        for i2 := 0; i2 < len(ruckSacks[i]); i2++ {
            priorities[i][i2] = priority(ruckSacks[i][i2])
        }
    }
    // find common priority per group
    for i := 0; i < len(priorities); i+=3 {
        for _, v := range priorities[i] {
            if contains(priorities[i+1], v) && contains(priorities[i+2], v) {
                priority_sum += v
                break
            }
        }
    }

    return fmt.Sprint(priority_sum)
}

func parse(input_file string) [][]byte {
    file, _ := os.Open(input_file)
    reader := bufio.NewReader(file)


    ruckSack := make([][]byte, 0)
    var err error = nil
    var line []byte
    for err == nil {
        line, err = reader.ReadBytes('\n')
        if len(line) > 0 {
            ruckSack = append(ruckSack, line[:len(line)-1])
        }
    }
    return ruckSack
}
