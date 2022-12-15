package day13

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type List struct {
    elems []Elem
}

type Value = int

type Elem struct {
    list *List
    val Value
}

func (e Elem) String() string {
    if e.list == nil {
        return fmt.Sprint(e.val)
    }else {
        ss := make([]string, 0)
        for _, e := range e.list.elems { ss = append(ss, e.String()) }
        return "[" + strings.Join(ss, ",") + "]"
    }
}

// Create type that supports for sort interface
type ElemList []Elem

func (x ElemList) Len() int {
    return len(x)
}

func (x ElemList) Less(i,j int) bool {
    return cmp(x[i], x[j]) == 1
}

func (x ElemList) Swap(i,j int) {
    tmp := x[i]
    x[i] = x[j]
    x[j] = tmp
}


const (
    NDet int = iota
    Correct
    Wrong
)

func Part1(input_file string) string {
    r := parse(input_file)
    answer := 0
    for i := 0; i < len(r) / 2; i++ {
        e1 := r[i*2]
        e2 := r[i*2+1]
        c := cmp(e1, e2)
        if c == Correct { answer += i + 1 }
    }
    return fmt.Sprint(answer)
}

func Part2(input_file string) string {
    e1 := Elem{&List{[]Elem{{&List{[]Elem{{nil, 2}}}, 0}}}, 0}
    e2 := Elem{&List{[]Elem{{&List{[]Elem{{nil, 6}}}, 0}}}, 0}
    r := ElemList(parse(input_file))
    r = append(r, e1, e2)
    sort.Sort(r)
    answer := 1
    for i := range r {
        if r[i] == e1 || r[i] == e2 {
            answer *= i + 1
        }
    }
    return fmt.Sprint(answer)
}

func cmp(elem1, elem2 Elem) int {
    e1 := elem1
    e2 := elem2
    if e1.list == nil && e2.list == nil { // compare two values
        if e1.val < e2.val { return Correct }
        if e1.val == e2.val { return NDet }
        return Wrong
    }

    // Convert one to list
    if e1.list == nil {
        elems := make([]Elem, 1)
        elems[0] = Elem{nil, e1.val}
        e1 = Elem{&List{elems}, 0}
    }

    if e2.list == nil {
        elems := make([]Elem, 1)
        elems[0] = Elem{nil, e2.val}
        e2 = Elem{&List{elems}, 0}
    }

    i := 0
    for i < min(len(e1.list.elems), len(e2.list.elems)) {
        c := cmp(e1.list.elems[i], e2.list.elems[i])
        if c == Correct || c == Wrong { return c }
        i++
    }
    // The inner elements did not result in a conclusion
    if len(e1.list.elems) == len(e2.list.elems) { return NDet }
    if i >= len(e2.list.elems) { return Wrong }
    return Correct
}


// -- PARSING ---------
type ParseRes struct {
    elem Elem
    rem string
}

func parse(input_file string) []Elem {
    f, _ := os.Open(input_file)
    b := bufio.NewScanner(f)
    elems := make([]Elem, 0)
    for b.Scan() {
        l1 := b.Text()
        e1 := parseElem(l1).elem
        b.Scan()
        l2 := b.Text()
        e2 := parseElem(l2).elem
        elems = append(elems, e1, e2)
        b.Scan() // get rid of redundant newline
    }

    return elems
}

func parseElem(ss string) ParseRes {
    if ss[0] == '[' {
        return parseList(ss)
    }else {
        return parseValue(ss)
    }
}

func parseList(ss string) ParseRes {
    elems := make([]Elem, 0)
    s := ss[1:]
    // Parse one element at a time (no backtracking)
    for {
        if s[0] == ']' {break}
        if s[0] == ',' {s = s[1:]} // ignore optional comma
        res := parseElem(s)
        s = res.rem
        elems = append(elems, res.elem)
    }

    l := List{elems}
    return ParseRes{
        Elem {&l, 0},
        s[1:],
    }
}

func parseValue(ss string) ParseRes {
    // Find index of where value ends
    i1 := ic(strings.Index(ss, "["))
    i2 := ic(strings.Index(ss, ","))
    i3 := ic(strings.Index(ss, "]"))
    i := min(i1,i2)
    i = min(i, i3)

    v, e := strconv.Atoi(string(ss[:i]))
    if e != nil {panic(fmt.Sprint(e))}
    return ParseRes{
        Elem {nil, v},
        ss[i:],
    }
}

func min(a,b int) int {
    if a < b {return a}
    return b
}

func max(a,b int) int {
    if a > b {return a}
    return b
}

func ic(a int) int {
    if a == -1 {return int(^uint(0) >> 1)}
    return a
}
