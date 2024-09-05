package main

import (
	"fmt"
	"os"
	"slices"
	"strings"
)

func main() {
	urlsList := strings.Split(readFile("config/lists.txt"), "\r\n")
	manual := parseList(readFile("config/manual.txt"))

	// Add lists
	allSet := make(map[string]struct{}, len(manual))
	addAllToSet(manual, allSet)
	for _, url := range urlsList {
		if url == "" || url[0] == '#' {
			continue
		}

		addAllToSet(getList(url), allSet)
	}

	suffixes := parseList(readFile("config/suffixes.txt"))

	// Move from a set to a slice
	i := 0
	all := make([]string, len(allSet))
	for domain := range allSet {
		if isExtempted(domain, suffixes) {
			continue
		}

		all[i] = domain
		i++
	}
	all = append(all, suffixes...)

	// Filter list and write to file
	os.MkdirAll("output/chunks", os.ModeDir)
	all = FilterDomains(all).Get()
	saveListToFile(all, "output/all.txt")

	// Save chunks
	i = 0
	for chunk := range slices.Chunk(all, 1000) {
		saveListToFile(chunk, fmt.Sprintf("output/chunks/%d.csv", i))
		i++
	}
}

func addAllToSet(l []string, m map[string]struct{}) {
	for _, v := range l {
		m[v] = struct{}{}
	}
}

func readFile(fileName string) string {
	data, err := os.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	return string(data)
}
