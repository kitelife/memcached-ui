package main

import (
		"os"
		"log"
		"fmt"
		"bufio"
		"strings"
		"strconv"
		// "unicode/utf8"
)

func main() {
	file, err := os.Open("serialize.txt")
	if err != nil {
		log.Fatal(err)
	}

	reader := bufio.NewReader(file)

	fmt.Printf("%q", parse(reader))
}
