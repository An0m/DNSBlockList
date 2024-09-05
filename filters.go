package main

import (
	"slices"
	"strings"
)

type Filter struct {
	SubDomains  []string
	SecondLevel []string
	parents     []string
}

func FilterDomains(domains []string) *Filter {
	f := Filter{}
	slices.Sort(domains)

	for _, domain := range domains {
		if strings.Count(domain, ".") < 2 {
			f.SecondLevel = append(f.SecondLevel, domain)
		} else {
			f.SubDomains = append(f.SubDomains, domain)

			// Get parent domain
			notTopLevel := domain[:strings.LastIndex(domain, ".")]
			f.parents = append(f.parents, domain[strings.LastIndex(notTopLevel, ".")+1:])
		}
	}

	return &f
}

func (f *Filter) getIndexes() []int {
	indexes := []int{}
	for i, domain := range f.parents {
		_, found := slices.BinarySearch(f.SecondLevel, domain) // The second level domains are sorted, the parents are not
		if !found {
			continue
		}

		indexes = append(indexes, i)
	}

	return indexes
}

func (f *Filter) filter(indexes []int) []string {
	finalDomains := make([]string, 0, len(f.SecondLevel)+len(f.SubDomains)-len(indexes))

	for i := 0; i < len(f.SubDomains); i++ {
		if _, found := slices.BinarySearch(indexes, i); found { // The indexes are thus also sorted
			continue
		}

		finalDomains = append(finalDomains, f.SubDomains[i])
	}

	return finalDomains
}

func (f *Filter) Get() []string {
	finalDomains := f.filter(f.getIndexes())
	slices.Sort(finalDomains)
	return finalDomains
}
