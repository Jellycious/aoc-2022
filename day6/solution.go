package day6

import (
	"fmt"
	"os"
)

func Part1(input_file string) string {
    bs, _ := os.ReadFile(input_file)
    s := string(bs)
    return distinct(s, 4)
}

func Part2(input_file string) string {
    bs, _ := os.ReadFile(input_file)
    s := string(bs)
    return distinct(s, 14)
}

// Checks whether a duplicate element is present
func dup(s string) bool {

    for i := range s {
        for j := i+1; j < len(s); j++ {
            if s[i] == s[j] {return true}
        }
    }

    return false
}

func distinct(s string, l int) string {
    for i := 0; i < len(s)-l; i++ {
        if !dup(s[i:i+l]) {
            return fmt.Sprint(i+l)
        }
    }

    panic("Shittttt")
}

