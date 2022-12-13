package day2

import (
	"bufio"
	"fmt"
	"os"
)

//-- DATA STRUCTURES ---------------------------------------
type HandShape int
const (
    A HandShape = iota  // Rock
    B                   // Paper
    C                   // Scissor
)

type Response int
const (
    X Response = iota   // Rock
    Y                   // Paper
    Z                   // Scissor
)

type Result int
const (
    Lose Result = iota
    Draw
    Win
)

type Round struct {
    handshape HandShape
    response Response
}


//-- UTILITY -----------------------------------------------

// Returns Result of round played with provided HandShape and Response
func roundResult(h HandShape, r Response) Result {
    if (int(h) + 1) % 3 == int(r) {
        // You win
        return Win
    }else if (int(r) + 1) % 3 == int(h) {
        // You lose
        return Lose
    }else {
        // You draw
        return Draw
    }
}

//-- SOLUTIONS ---------------------------------------------

// Solution of Part 1
func Part1(input_file string) string {
    rounds := parse(input_file)
    score := 0
    score_map := map[Result]int {
        Win: 6,
        Draw: 3,
        Lose: 0,
    }

    // Loop through rounds and increment score
    for _, r := range rounds {
        score += int(r.response) + 1
        score += score_map[roundResult(r.handshape, r.response)]
    }
    return fmt.Sprint(score)
}

// Solution of Part 2
func Part2(input_file string) string {
    rounds := parse(input_file)
    score := 0
    score_map := map[Result]int {
        Win: 6,
        Draw: 3,
        Lose: 0,
    }
    response := []Response {X,Y,Z}

    for _, r := range rounds {
        shape := r.handshape

        // Try all possible responses to get desired result
        // More efficient would be to construct a predefined map
        for _, cand_r := range response {
            res := roundResult(shape, cand_r)

            // Check if res is equal to desired result
            if res == Result(int(r.response)) {
                score += score_map[res]
                score += int(cand_r) + 1
            }
        }

    }

    return fmt.Sprint(score)
}

//-- PASING ------------------------------------------------
func parse(input_file string) []Round {
    file, _ := os.Open(input_file)
    reader := bufio.NewReader(file)
    hsMap := map[uint8]HandShape {
        'A': A,
        'B': B,
        'C': C,
    }
    rsMap := map[uint8]Response {
        'X': X,
        'Y': Y,
        'Z': Z,
    }
    rounds := make([]Round, 0)
    for {
        line, err := reader.ReadString('\n')
        if err != nil {break}
        round := Round {hsMap[line[0]], rsMap[line[2]]}
        rounds = append(rounds, round)
    }
    return rounds
}
