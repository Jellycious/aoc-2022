package day1

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
    "sort"
)

func checkErr(e error) {
    if e != nil {
        panic(e)
    }
}

func sumList(l []int) int {
    sum := 0
    for _, v := range l{
        sum += v
    }
    return sum
}

func Part1(input_file string) string {
    cals := parse(input_file)
    sums := make([]int, len(cals))

    // Sum every sublist
    for _, elf := range cals {
        sums = append(sums, sumList(elf))
    }

    // find maximum
    max := 0
    for _, sum := range sums {
        if sum > max {max = sum}
    }

    return fmt.Sprint(max)
}

func Part2(input_file string) string {
    cals := parse(input_file)
    sums := make([]int, len(cals))
    // Sum every sublist
    for _, elf := range cals {
        sum := 0
        for _, cal := range elf {
            sum += cal
        }
        sums = append(sums, sum)
    }

    // Sort list
    sort.Sort(sort.Reverse(sort.IntSlice(sums)))
    top_three := sumList(sums[:3])
    return fmt.Sprint(top_three)
}


func parse(input_file string) [][]int {
    file, err := os.Open(input_file)
    checkErr(err)
    reader := bufio.NewReader(file)

    calories := make([][]int, 0)
    elf_calories := make([]int, 0)

    // PARSING
    for {
        line, err := reader.ReadString('\n')

        if err != nil {
            calories = append(calories, elf_calories)
            break
        }

        if line == "\n" {
            calories = append(calories, elf_calories)
            elf_calories = make([]int, 0)
            continue
        }
        val, err := strconv.Atoi(line[:len(line)-1])
        if err != nil {panic(err)}
        elf_calories = append(elf_calories, val)
    }

    return calories
}
