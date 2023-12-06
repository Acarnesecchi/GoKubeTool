package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {

	// Randomly choose a file to open
	fileName := "failure.log"
	if rand.Intn(2) == 0 { // Randomly select between 0 and 1
		fileName = "success.log"
	}

	// Open the chosen file
	file, err := os.Open(fileName)
	if err != nil {
		check(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
		random := time.Duration(rand.Intn(10))
		time.Sleep(random * time.Second)
	}

	check(scanner.Err())

	// Set exit code based on file opened
	if fileName == "success.log" {
		os.Exit(0)
	} else {
		os.Exit(1)
	}
}
