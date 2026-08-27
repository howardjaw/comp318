package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

func main() {
	file_name := flag.String("file_name", "", "select file")
	flag.Parse()
	file, err := os.Open(*file_name)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	var count int

	for scanner.Scan() {
		count++
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("scanner error: ", err)
	}

	fmt.Println("number of lines: ", count)
}
