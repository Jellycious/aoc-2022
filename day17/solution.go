package day17

import (
	"fmt"
	"os"
)

type Pos struct {
    x int
    y int
}

type Rock []Pos

type Jet []byte

func Part1(input_file string) string {
    return fmt.Sprint(parse(input_file))
}

func parse(input_file string) Jet {
    bs, _ := os.ReadFile(input_file)
    jet := Jet(make([]byte, 0))

    for _, b := range bs[:len(bs)-1] {
        jet = append(jet, b)
    }

    return jet
}
