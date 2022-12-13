package day5

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Move struct {
    quantity int
    start int
    end int
}

type Stack []rune

type Crates struct {
    stacks []Stack
}

type State struct {
    crates Crates
    moves []Move
}

func (p Stack) String() string {
    return string(p)
}

func (p Crates) String() string {
    stacks := make([]string, len(p.stacks))
    for i := range p.stacks {
        stacks[i] = strings.Join([]string{fmt.Sprint(i), p.stacks[i].String()}, ": ")
    }
    return strings.Join(stacks, "\n")
}

func Part1(input_file string) string {
    state := parse(input_file)
    crates := state.crates
    instructions := state.moves

    for _, move := range instructions {
        picked := crates.stacks[move.start-1][:move.quantity]
        // Drop values
        crates.stacks[move.start-1] = crates.stacks[move.start-1][move.quantity:]

        // Extend other stack with new crates
        new_stack := make(Stack, len(crates.stacks[move.end-1])+len(picked))
        for i := range picked {
            new_stack[i] = picked[len(picked)-1-i]
        }

        for i, c := range crates.stacks[move.end-1] {
            new_stack[i+move.quantity] = c
        }

        crates.stacks[move.end-1] = new_stack

    }

    top_crates := make([]rune, len(crates.stacks))
    for i := range crates.stacks {
        top_crates[i] = crates.stacks[i][0]
    }


    return fmt.Sprint(string(top_crates))
}

func Part2(input_file string) string {
    state := parse(input_file)
    crates := state.crates
    instructions := state.moves

    for _, move := range instructions {
        picked := crates.stacks[move.start-1][:move.quantity]
        crates.stacks[move.start-1] = crates.stacks[move.start-1][move.quantity:]

        new_stack := make(Stack, len(crates.stacks[move.end-1])+len(picked))
        for i, p := range picked {
            // The only change compared to Part1
            new_stack[i] = p
        }

        for i, c := range crates.stacks[move.end-1] {
            new_stack[i+move.quantity] = c
        }

        crates.stacks[move.end-1] = new_stack

    }

    top_crates := make([]rune, len(crates.stacks))
    for i := range crates.stacks {
        top_crates[i] = crates.stacks[i][0]
    }


    return fmt.Sprint(string(top_crates))
}

// -- PARSING ----------------------------------------------
func parse_crates(stack_string string) Crates {
    lines := strings.Split(stack_string, "\n")
    indices := strings.Split(strings.TrimSpace(lines[len(lines)-1]), " ")
    stack_count, _ := strconv.Atoi(indices[len(indices)-1])

    // Create stacks object
    crates := Crates{
        stacks: make([]Stack, stack_count),
    }
    // initialize arrays inside of stacks
    for i := 0; i < stack_count; i++ {
        crates.stacks[i] = make(Stack, 0)
    }

    for _, srow := range lines[:len(lines)-1] {
        // get rid of '[' ']' and ' '

        for i := 0; i < stack_count; i++ {
            index := 4*i + 1
            c := rune(srow[index])
            if c != ' ' {
                crates.stacks[i] = append(crates.stacks[i], c)
            }
        }
    }

    return crates
}

func parse_moves(moves_string string) []Move {
    moves := make([]Move, 0)

    lines := strings.Split(strings.TrimSpace(moves_string), "\n")

    for _, l := range lines {
        s := strings.ReplaceAll(l, "move ", "")
        s = strings.ReplaceAll(s, " from ", " ")
        s = strings.ReplaceAll(s, " to ", " ")
        numbers := strings.Split(s, " ")
        n1, _ := strconv.Atoi(numbers[0])
        n2, _ := strconv.Atoi(numbers[1])
        n3, _ := strconv.Atoi(numbers[2])

        moves = append(moves, Move {
            quantity: n1,
            start: n2,
            end: n3,

        })


    }

    return moves
}

func parse(input_file string) State {
    sb, _ := os.ReadFile(input_file)
    s := string(sb)

    split := strings.Split(s, "\n\n")
    crates := parse_crates(split[0])
    moves := parse_moves(split[1])
    // We can find out the number of crates by reading last integer

    return State{
        crates: crates,
        moves: moves,
    }
}
