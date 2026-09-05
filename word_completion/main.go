package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/howardjaw/comp318/word_completion/base"
	"github.com/howardjaw/comp318/word_completion/frequency"
)

type Completer interface {
	Complete(string) []string
}

func main() {
	// Input flag to specify dictionary file, which must consists of lines that
	// are a string, followed by a space, followed by a number.
	dictFlag := flag.String("d", "", "dictionary file name")
	compFlag := flag.String("c", "base", "completer type [base, frequency]")

	flag.Parse()

	// Ensure that a filename was specified on the command line.
	if *dictFlag == "" {
		fmt.Printf("Must specify a dictionary file with the -d flag\n")
		return
	}

	// Open dictionary file.
	file, err := os.Open(*dictFlag)
	if err != nil {
		fmt.Printf("Unable to open file %s: %v\n", *dictFlag, err)
		return
	}
	defer file.Close()

	// Read dictionary file.
	dict := make(map[string]int)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, " ")
		if len(fields) != 2 {
			fmt.Printf("Invalid line: %s\n", line)
			return
		}
		word := fields[0]
		count, err := strconv.Atoi(fields[1])
		if err != nil {
			fmt.Printf("Invalid count: %s\n", fields[1])
			return
		}

		_, ok := dict[word]
		if ok {
			fmt.Printf("Duplicate word: %s\n", word)
			return
		}

		dict[word] = count
	}

	err = scanner.Err()
	if err != nil {
		fmt.Printf("Unable to read file %s: %v\n", *dictFlag, err)
		return
	}

	// Instantiate correct Completer
	var completer Completer

	if *compFlag == "base" {
		completer = base.New(dict)
	} else if *compFlag == "frequency" {
		completer = frequency.New(dict)
	} else {
		fmt.Printf("Invalid completer type: %s\n", *compFlag)
		return
	}

	// Read input from user and provide completions.
	// Ctrl-C to exit.
	for {
		var prefix string

		fmt.Printf("String: ")
		fmt.Scanln(&prefix)
		if prefix == "" {
			fmt.Println("Please type one or more characters")
			continue
		}
		fmt.Println(completer.Complete(prefix))
	}
}
