package main

import (
	_ "embed"
	"encoding/xml"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"text/template"
)

const srcURL = "https://www.six-group.com/dam/download/financial-information/data-center/iso-currrency/lists/list-one.xml"

type ISO4217 struct {
	Table struct {
		Entries []Entry `xml:"CcyNtry"`
	} `xml:"CcyTbl"`
}

type Entry struct {
	Numeric string `xml:"CcyNbr"`
	Alpha   string `xml:"Ccy"`
}

//go:embed data.tmpl
var tmpl string

func main() {
	// load XML data
	resp, err := http.Get(srcURL)
	if err != nil {
		log.Fatalf("failed to fetch XML: %v", err)
	}
	defer resp.Body.Close()

	// parse data
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("failed to read response body: %v", err)
	}

	var doc ISO4217
	if err := xml.Unmarshal(data, &doc); err != nil {
		log.Fatalf("failed to unmarshal XML: %v", err)
	}

	// filter, deduplicate and sort data
	var entries []Entry
	seen := make(map[string]bool)
	for _, e := range doc.Table.Entries {
		if e.Numeric != "" && e.Alpha != "" {
			if seen[e.Numeric] {
				continue
			}
			seen[e.Numeric] = true
			entries = append(entries, e)
		}
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Numeric < entries[j].Numeric
	})

	// generate output file
	f, err := os.Create("utils/iso4217_code.go")
	if err != nil {
		log.Fatalf("failed to create output file: %v", err)
	}
	defer f.Close()

	t := template.Must(template.New("iso4217").Parse(tmpl))
	if err := t.Execute(f, entries); err != nil {
		log.Fatalf("failed to execute template: %v", err)
	}
}
