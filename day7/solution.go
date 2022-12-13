package day7

import (
	"bufio"
	"fmt"
	"os"

	"github.com/alecthomas/participle/v2"
)

type Directory struct {
    Parent *Directory
    Name string
    SubDirs map[string]Directory
    Files map[string]int
}

func (d Directory) String() string {
    s := fmt.Sprintf("Directory %s:\n", d.Name)

    for _, dir := range d.SubDirs {
        s = s + "   dir " + dir.Name + "\n"
    }
    for filename, size := range d.Files {
        s = s + "   " + filename +": "+fmt.Sprintln(size)
    }
    return s
}

func createDir(parent *Directory, name string) Directory {
    return Directory {
        parent,
        name,
        make(map[string]Directory, 0),
        make(map[string]int, 0),
    }
}

func Part1(input_file string) string {
    lines := parse(input_file)
    dir := constructFS(lines)

    tsizes := make(map[string]int, 0)
    DFS(dir, "/", &tsizes)

    answer := 0
    for _, val := range tsizes {
        if val <= 100000 { answer += val }
    }

    return fmt.Sprint(answer)
}

func Part2(input_file string) string {
    lines := parse(input_file)
    dir := constructFS(lines)

    tsizes := make(map[string]int, 0)
    DFS(dir, "/", &tsizes)

    answer := 70000000

    for dir, size := range tsizes {
        if 70000000 - tsizes["/"] + size < 30000000 {delete(tsizes, dir)}
    }

    for _, size := range tsizes {
        if size < answer {
            answer = size
        }
    }

    return fmt.Sprint(answer)
}

func DFS(d Directory, abspath string, tsizes *map[string]int) {
    dtsize := 0
    for name, dir := range(d.SubDirs) {
        new_path := abspath+name+"/"
        DFS(dir, new_path, tsizes)
        size := (*tsizes)[new_path]
        dtsize += size
    }

    tsize := 0
    for _, size := range(d.Files) {
        tsize += size
    }

    (*tsizes)[abspath] = tsize + dtsize
}



func constructFS(lines []Line) Directory {

    cur_dir := &Directory {
        nil,
        "/",
        make(map[string]Directory, 0),
        make(map[string]int, 0),
    }

    for _, line := range lines[1:] {
        // Check line
        if line.isCmd {
            if line.Cmd.Cmd == ls {continue}

            if line.Cmd.Arg == ".." {
                cur_dir = cur_dir.Parent
                continue
            }
            // Check whether dir is present
            new_dir := createDir(cur_dir, line.Cmd.Arg)
            cur_dir.SubDirs[line.Cmd.Arg] = new_dir
            cur_dir = &new_dir

        }else {
            // Check files
            t := fmt.Sprintf("%T", line.Out)
            if t == "day7.File" {
                file := line.Out.(File)
                cur_dir.Files[file.Name] = file.Size
            }
        }

    }

    for cur_dir.Parent != nil {
        cur_dir = cur_dir.Parent
    }

    return *cur_dir
}


// -- PARSING --------------------
type Cmd string
const (
    cd Cmd = "cd"
    ls = "ls"
)

type Line struct {
    Cmd Command
    Out Output
    isCmd bool
}

type Command struct {
    Cmd Cmd `"$" @( "cd" | "ls" )`
    Arg string `(@("."".") | @"/" | @Ident)?`
}

type Output interface {output()}

type File struct {
    Size int `@Int`
    Name string `(@Ident @"." @Ident | @Ident)`
}
func (File) output() {}


type Dir struct {
    Name string `"dir" @Ident`
}
func (Dir) output() {}

func parse(input_string string) []Line {
    f, _ := os.Open(input_string)
    scanner := bufio.NewScanner(f)

    cmd_parser := participle.MustBuild[Command]()

    out_parser := participle.MustBuild[Output](
        participle.Union[Output](File{}, Dir{}),
    )

    lines := make([]Line, 0)

    for scanner.Scan() {
        l := scanner.Text()

        lp := Line{}

        // Check whether command or output
        if l[0] == '$' {
            // Command
            cmd, _ := cmd_parser.ParseString("", l)
            lp.Cmd = *cmd
            lp.isCmd = true

        }else {
            // Output
            out, _ := out_parser.ParseString("", l)
            lp.Out = *out
            lp.isCmd = false
        }

        lines = append(lines, lp)
    }

    return lines
}
