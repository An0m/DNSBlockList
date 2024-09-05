package main

import (
	"fmt"
	"os"
	"slices"
	"strings"
	"time"
)

type UnorderedList struct {
	slice   []string
	removed int
}

func (l *UnorderedList) Add(items ...string) {
	l.slice = append(l.slice, items...)
}
func (l *UnorderedList) getLastIndex() int {
	return len(l.slice) - l.removed
}
func (l *UnorderedList) Remove(index int) {
	l.slice[index] = l.slice[l.getLastIndex()-1]
}
func (l *UnorderedList) Get() []string {
	return l.slice[:l.getLastIndex()]
}

func main() {
	data, err := os.ReadFile("output/all.txt")
	if err != nil {
		panic(err)
	}

	domains := strings.Split(string(data), "\r\n")
	subdomains := []string{}
	parents := []string{}
	secondlevel := []string{}
	indexes := []int{}
	startingTime := time.Now()

	fmt.Println("Preloading...")
	t := time.Now()
	for _, domain := range domains {
		if strings.Count(domain, ".") < 2 {
			secondlevel = append(secondlevel, domain)
		} else {
			subdomains = append(subdomains, domain)

			// Get parent domain
			notTopLevel := domain[:strings.LastIndex(domain, ".")]
			parents = append(parents, domain[strings.LastIndex(notTopLevel, ".")+1:])
		}
	}
	fmt.Printf("Done! %v\n", time.Since(t))

	fmt.Println("Searching indexes...")
	t = time.Now()
	for i, domain := range parents {
		_, found := slices.BinarySearch(secondlevel, domain) // The second level domains are sorted, the parents are not
		if !found {
			continue
		}

		indexes = append(indexes, i)
	}
	fmt.Printf("Done! %v\n", time.Since(t))

	fmt.Printf("Adding items to slice! %v\n", time.Since(t))
	finalDomains := secondlevel
	t = time.Now()
	for i := 0; i < len(subdomains); i++ {
		if _, found := slices.BinarySearch(indexes, i); found { // The indexes are thus also sorted
			continue
		}

		finalDomains = append(finalDomains, subdomains[i])
	}
	slices.Sort(finalDomains)
	fmt.Printf("Done! %v\n", time.Since(t))

	fmt.Println("All done! Writing to file")
	os.WriteFile("output/filtered.txt", []byte(strings.Join(finalDomains, "\n")), 0)
	fmt.Printf(
		"Total time elapsed: %v | Final domain number: %d/%d (%v%% less)",
		time.Since(startingTime),
		len(finalDomains), len(domains),
		100-(float32(len(finalDomains))/float32(len(domains)))*100,
	)

}
