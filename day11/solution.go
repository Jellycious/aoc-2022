package day11

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/Jellycious/aoc2022/utils"
)

const (
    Mul = '*'
    Add = '+'
)

type Op struct {
    opType rune
    is_old bool
    val int
}

type Monkey struct {
    items []int
    op Op
    divBy int
    ifTrue int
    ifFalse int
}

var MAGIC_NUMBER int

func Part1(input_file string) string {
    return solve(input_file, true)
}

func Part2(input_file string) string {
    return solve(input_file, false)
}

func solve(input_file string, part1 bool) string {
    monkeys := parse(input_file)
    inspections := make(map[int]int, len(monkeys))

    //compute magic number
    MAGIC_NUMBER = 1

    for i := 0; i < len(monkeys); i++ {
        inspections[i] = 0
        MAGIC_NUMBER *= monkeys[i].divBy
    }

    rounds := 10000
    if part1 {rounds = 20}

    for i := 0; i < rounds; i++ {
        if part1 {
            simulateRound1(&monkeys, &inspections)
        }else {
            simulateRound2(&monkeys, &inspections)
        }
    }

    // Put counts in list and sort
    c := make([]int, len(inspections))
    i := 0
    for _, v := range inspections {
        c[i] = v
        i++
    }
    sort.Ints(c)

    answer := c[len(c)-1] * c[len(c)-2]

    return fmt.Sprint(answer)
}

func simulateRound1(monkeys *[]Monkey, inspections *map[int]int) {
    for m := range (*monkeys) {
        monkey := (*monkeys)[m]

        for _, w := range monkey.items {
            // Compute new worry level
            (*inspections)[m] = (*inspections)[m] + 1
            worry_level := w
            operand := monkey.op.val
            if monkey.op.is_old {operand = worry_level}

            if monkey.op.opType == Mul {
                worry_level = (worry_level * operand) / 3
            }else if monkey.op.opType == Add {
                worry_level = (worry_level + operand) / 3
            }

            // Throw to other monkey condition
            var otherItems *[]int
            if worry_level % monkey.divBy == 0 {
                otherItems = &(*monkeys)[monkey.ifTrue].items
            }else {
                otherItems = &(*monkeys)[monkey.ifFalse].items
            }
            r := append(*otherItems, worry_level)
            *otherItems = r
        }
        // Remove all items in monkeys list
        (*monkeys)[m].items = make([]int, 0)
    }
}


func simulateRound2(monkeys *[]Monkey, inspections *map[int]int) {
    for m := range (*monkeys) {
        monkey := (*monkeys)[m]

        for _, w := range monkey.items {
            // Compute new worry level
            (*inspections)[m] = (*inspections)[m] + 1
            worry_level := w
            operand := monkey.op.val
            if monkey.op.is_old {operand = worry_level}

            if monkey.op.opType == Mul {
                worry_level = (worry_level * operand) % MAGIC_NUMBER
            }else if monkey.op.opType == Add {
                worry_level = (worry_level + operand) % MAGIC_NUMBER
            }

            // Throw to other monkey condition
            var otherItems *[]int
            if worry_level % monkey.divBy == 0 {
                otherItems = &(*monkeys)[monkey.ifTrue].items
            }else {
                otherItems = &(*monkeys)[monkey.ifFalse].items
            }
            r := append(*otherItems, worry_level)
            *otherItems = r
        }
        // Remove all items in monkeys list
        (*monkeys)[m].items = make([]int, 0)
    }
}


func parse(input_file string) []Monkey {
    f, _ := os.Open(input_file)
    scanner := bufio.NewScanner(f)
    scanner.Split(utils.ScanSubstr([]byte("\n\n")))

    monkeys := make([]Monkey, 0)

    for scanner.Scan() {
        monkey := strings.Split(scanner.Text(), "\n")
        // Parse items
        items_s := strings.Split(string(monkey[1][18:]), ", ")
        items := make([]int, len(items_s))
        for i, s := range items_s {
            v, _ := strconv.Atoi(s)
            items[i] = v
        }

        // Parse Operation
        op := Op{}
        op.opType = rune(monkey[2][23])
        op.is_old = monkey[2][25] == 'o'
        v, _ := strconv.Atoi(string(monkey[2][25:]))
        op.val = v

        // Parse divisible by
        div_by,_ := strconv.Atoi(monkey[3][21:])

        // Parse ifTrue 
        if_true,_ := strconv.Atoi(monkey[4][29:])
        if_false,_ := strconv.Atoi(monkey[5][30:])

        monkeys = append(monkeys, Monkey {items, op, div_by, if_true, if_false})
    }

    return monkeys
}
