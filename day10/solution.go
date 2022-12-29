package day10

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

type Instr struct {
    is_noop bool
    val int
}

func Part1(input_file string) string {
    instrs := parse(input_file)
    cycle := 0
    x := 1
    answer := 0
    for _, instr := range instrs {
        cycle++ // Begin new cycle
        answer += checkSignal(cycle, x)

        if instr.is_noop {
            continue
        }

        cycle++
        answer += checkSignal(cycle, x)
        x += instr.val // update register
    }
    return fmt.Sprint(answer)
}

func Part2(input_file string) string {
    instrs := parse(input_file)
    cycle := 0
    x := 1
    crt := ""
    for _, instr := range instrs {
        cycle++ // Begin new cycle
        //fmt.Printf("cycle: %d x: %d\n", cycle, x)
        crt += draw(x, cycle)
        //fmt.Println(crt)
        //fmt.Println()

        if instr.is_noop {
            continue
        }

        cycle++
        //fmt.Printf("cycle: %d x: %d\n", cycle, x)
        crt += draw(x, cycle)
        //fmt.Println(crt)
        //fmt.Println()
        x += instr.val // update register
    }
    return crt
}

func checkSignal(cycle int, x int) int {
    if (cycle - 20) % 40 == 0 {
        return x * cycle
    }
    return 0
}

func draw(x int, cycle int) string {
    var suffix string = ""
    if (cycle) % 40 == 0 && cycle != 1 {suffix = "\n"}
    c := (cycle - 1) % 40
    if c >= (x-1) && c <= (x+1) {return "X"+suffix}
    return " "+suffix
}

func parse(input_file string) []Instr {
    f, _ := os.Open(input_file)
    scanner := bufio.NewScanner(f)

    instructions := make([]Instr, 0)
    for scanner.Scan() {
        s := scanner.Text()
        if s[0] == 'n' {
            instructions = append(instructions, Instr{true, 0})
        }else {
            n, _ := strconv.Atoi(string(s[5:]))
            instructions = append(instructions, Instr{false, n})
        }
    }
    return instructions
}
