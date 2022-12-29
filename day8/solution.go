package day8

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
    "sort"
)

func Part1(input_file string) string {
    trees := parse(input_file)
    trees_t := transpose(trees)

    visible_count := len(trees)*2 + len(trees[0])*2 - 4
    for ri := 1; ri < len(trees)-1; ri++ {
        for ci := 1; ci < len(trees[0])-1; ci++ {
            // Check whether the tree is visible
            t := trees[ri][ci]
            left := trees[ri][:ci]
            right := trees[ri][ci+1:]
            top := trees_t[ci][:ri]
            bottom := trees_t[ci][ri+1:]
            is_visible := visible(t, left) ||visible(t, right) || visible(t, top) || visible(t, bottom)
            if is_visible {visible_count++}
        }
    }
    return fmt.Sprint(visible_count)
}

func Part2(input_file string) string {
    trees := parse(input_file)
    trees_t := transpose(trees)

    max_visible_score := 0
    for ri := 0; ri < len(trees); ri++ {
        for ci := 0; ci < len(trees[0]); ci++ {
            // Check whether the tree is visible
            t := trees[ri][ci]
            left := visiblec(t, trees[ri][:ci], true)
            right := visiblec(t, trees[ri][ci+1:], false)
            top := visiblec(t, trees_t[ci][:ri], true)
            bottom := visiblec(t, trees_t[ci][ri+1:], false)
            max_visible_score = max(max_visible_score, left * right * top * bottom)
        }
    }
    return fmt.Sprint(max_visible_score)
}

func ReverseSlice[T comparable](s []T) {
    sort.SliceStable(s, func(i, j int) bool {
        return i > j
    })
}

func max(a, b int) int {
    if a > b {
        return a
    }
    return b
}


func visible(t int, n []int) bool {
    visible := true
    for i := range n {
        visible = visible && (n[i] < t)
    }
    return visible
}

func visiblec(t int, n []int, reverse bool) int {
    visible := 0
    if reverse {
        for i := len(n) - 1; i >= 0; i-- {
            visible++
            if n[i] >= t {return visible}
        }
    }else {
        for i := range n {
            visible++
            if n[i] >= t {return visible}
        }
    }
    return visible
}

func transpose(m [][]int) [][]int {
    res := make([][]int, len(m[0]))
    for j := 0; j < len(m[0]); j++ {
        n := make([]int, len(m))
        for i := 0; i < len(m); i++ {
            n[i] = m[i][j]
        }
        res[j] = n
    }
    return res
}

func parse(input_file string) [][]int {
    f, _ := os.Open(input_file)
    scanner := bufio.NewScanner(f)

    rows := make([][]int, 0)
    for scanner.Scan() {
        row := scanner.Text()
        row_n := make([]int, len(row))

        for i, r := range row {
            v, _ := strconv.Atoi(string(r))
            row_n[i] = v
        }
        rows = append(rows, row_n)
    }
    return rows
}
