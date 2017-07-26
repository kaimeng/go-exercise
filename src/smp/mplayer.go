package main

import (
	"smp/library"
	"fmt"
	"strconv"
	"smp/mp"
	"bufio"
	"os"
	"strings"
)

var lib *library.MusicManager
var id int = 1
var ctr1, signal chan int

func handleLibCommand(tokens []string) {
	switch tokens[1] {
	case "list":
		for i := 0; i < lib.Len(); i++ {
			e, _ := lib.Get(i)
			fmt.Println(i+1, ":", e.Name, e.Artist, e.Sorce, e.Type)
		}
	case "add":
		if len(tokens) == 6 {
			id++
			lib.Add(&library.MusicEntry{strconv.Itoa(id),
										tokens[2], tokens[3], tokens[4], tokens[5]})
		} else {
			fmt.Println("USAGE: lib add <name><artist><source><type>")
		}
	case "remove":
		if len(tokens) == 3 {
			index, err := strconv.Atoi(tokens[2])
			if err != nil {
				lib.Remove(index)
				return
			}
		}
		fmt.Println("USAGE: lib remove <name>")
	default:
		fmt.Println("Unrecognized lib command:", tokens[1])
	}
}

func handlePlayCommand(tokens []string) {
	if len(tokens) != 2 {
		fmt.Println("USAGE: play <index>")
		return
	}

	index, err := strconv.Atoi(tokens[1])
	if err != nil {
		fmt.Println("USAGE: play <index>")
		return
	}
	e, err := lib.Get(index)
	if err != nil {
		fmt.Println("USAGE: play <index>")
		return
	}
	mp.Play(e.Sorce, e.Type)
}

func main() {
	fmt.Println(`
		Enter following commands to control the player:
		lib list -- View the existing music lib
		lib add <name><artist><source><type> -- Add a music to the music lib
		lib remove <name> -- Remove the specified music from the lib
		play <name> -- Play the specified music
		`)
	lib = library.NewMusicManager()
	r := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter command-> ")
		rawLine, _, _ := r.ReadLine()
		line := string(rawLine)
		if line == "q" || line == "e" {
			break
		}
		tokens := strings.Split(line, " ")
		if tokens[0] == "lib" {
			handleLibCommand(tokens)
		} else if tokens[0] == "play" {
			handlePlayCommand(tokens)
		} else {
			fmt.Println("Unrecognized command:", tokens[0])
		}
	}
}
