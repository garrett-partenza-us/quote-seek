package main

import (
	"fmt"
	"bufio"
	"os"
)

func GetEntries(c *Config, entries []string) []string {	

	file, err := os.Open(c.DatasetPath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		entries = append(entries, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return entries

}

func main() {

	config := GetConfig()
	var entries []string
	entries = GetEntries(&config, entries)
	fmt.Println(len(entries))

}
