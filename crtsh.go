package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"unicode"
)

const (
	RE_SUBDOMAIN_TPL = `^(?:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?\.)*%s$`
	RE_DOMAIN        = `^(?:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?\.)+[a-z0-9][a-z0-9-]{0,61}[a-z0-9]$`
)

type CrtshRecord struct {
	NameValue string `json:"name_value"`
}

func main() {
	useJson := flag.Bool("json", false, "Output in JSON format")
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		log.Fatal("Domain name is not specified")
	}
	if len(args) > 1 {
		log.Fatal("Invalid order or too many arguments")
	}

	domain := strings.ToLower(args[0])
	if ok, _ := regexp.MatchString(RE_DOMAIN, domain); !ok {
		log.Fatal("Invalid domain name")
	}

	data, err := fetch(domain)
	if err != nil {
		log.Fatal(err)
	}

	subdomainsMap, err := parse(data, domain)
	if err != nil {
		log.Fatal(err)
	}

	subdomains := sortedKeys(subdomainsMap)
	if *useJson {
		printJson(subdomains)
	} else {
		printText(subdomains)
	}
}

func fetch(domain string) ([]byte, error) {
	url := fmt.Sprintf("https://crt.sh/?q=%s&output=json", domain)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP Response: %s", resp.Status)
	}

	return io.ReadAll(resp.Body)
}

func parse(data []byte, domain string) (map[string]bool, error) {
	domains := make(map[string]bool)
	r := regexp.MustCompile(fmt.Sprintf(RE_SUBDOMAIN_TPL, domain))

	var records []CrtshRecord
	if err := json.Unmarshal(data, &records); err != nil {
		return nil, fmt.Errorf("JSON: %s", err)
	}

	for _, rec := range records {
		lines := strings.Split(strings.ToLower(rec.NameValue), "\n")
		for _, line := range lines {
			line = strings.TrimFunc(line, unicode.IsSpace)
			if r.MatchString(line) {
				domains[line] = true
			}
		}
	}

	return domains, nil
}

func printJson(subdomains []string) {
	bts, err := json.Marshal(subdomains)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(bts))
}

func printText(subdomains []string) {
	for _, subdomain := range subdomains {
		fmt.Println(subdomain)
	}
}

func sortedKeys(m map[string]bool) []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
