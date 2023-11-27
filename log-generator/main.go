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
	file, err := os.Open("sample.txt")
	if err != nil {
		check(err)
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
		random := time.Duration(rand.Intn(10))
		time.Sleep(random * time.Second)
	}

	check(scanner.Err())
	os.Exit(0)
}
