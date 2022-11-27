package main

import (
	"fmt"

	"github.com/kembo91/apc-test/bfs"
)

func main() {
	rt, err := bfs.NewRaceTracksFromFile("input.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, r := range rt {
		res := r.Race()
		if res == -1 {
			fmt.Println("No solution")
		} else {
			fmt.Printf("Optimal solution takes %v steps\n", res)
		}
	}
}
