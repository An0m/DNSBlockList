package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

var (
	DOMAIN_REGEX = regexp.MustCompile(`([a-z0-9]+(-[a-z0-9]+)*\.)+[a-z]{2,}`)
)

func parseList(data string) []string {
	data = strings.ReplaceAll(data, "\t", " ")
	data = strings.ReplaceAll(data, "\r", "")
	lines := strings.Split(data, "\n")

	domains := make([]string, 0, len(lines))
	for _, line := range lines {
		if len(line) < 4 || line[0] == '#' || line[0] == '!' {
			continue
		}

		domain := DOMAIN_REGEX.FindString(line)
		if domain == "" {
			continue
		}

		domain = strings.TrimPrefix(domain, "0.0.0.0")
		domain = strings.TrimPrefix(domain, "www.")

		//TODO: Implement filtering
		domains = append(domains, domain)
	}

	return domains
}

func getList(url string) []string {
	r, err := http.Get(url)
	if err != nil {
		log.Fatalf("Unable to resolve %s: %v", url, err)
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	return parseList(string(data))
}

func isExtempted(domain string, suffixes []string) bool {
	for _, suffix := range suffixes {
		if strings.HasSuffix(domain, suffix) {
			return true
		}
	}

	return false
}

func saveListToFile(list []string, filename string) {
	os.WriteFile(filename, []byte(strings.Join(list, "\n")), os.ModePerm)
}
