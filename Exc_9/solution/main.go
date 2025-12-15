package main

import (
	"bufio"
	"exc9/mapred"
	"fmt"
	"log"
	"os"
)

// Main function
func main() {
	// todo read file
	f, err := os.Open("res/meditations.txt")

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	var text []string

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		text = append(text, sc.Text())
	}
	if err := sc.Err(); err != nil {
		log.Fatal(err)
	}

	// todo run your mapreduce algorithm
	var mr mapred.MapReduce
	results := mr.Run(text)

	// todo print your result to stdout
	for w, c := range results {
		fmt.Println(w, c)
	}
}
