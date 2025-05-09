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

type Data struct {
	N2C []Entry
	C2N []Entry
}

//go:embed data.tmpl
var tmpl string

func buildN2C(entries []Entry) []Entry {
	var result []Entry
	seen := make(map[string]bool)
	for _, e := range entries {
		if seen[e.Numeric] {
			continue
		}
		seen[e.Numeric] = true
		result = append(result, e)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Numeric < result[j].Numeric
	})

	return result
}

func buildC2N(entries []Entry) []Entry {
	var result []Entry
	seen := make(map[string]bool)
	for _, e := range entries {
		if seen[e.Alpha] {
			continue
		}
		seen[e.Alpha] = true
		result = append(result, e)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Alpha < result[j].Alpha
	})

	return result
}

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

	// filter data
	var entries []Entry
	for _, e := range doc.Table.Entries {
		if e.Numeric != "" && e.Alpha != "" {
			entries = append(entries, e)
		}
	}

	n2c := buildN2C(entries)
	c2n := buildC2N(entries)

	// generate output file
	f, err := os.Create("utils/iso4217_data.go")
	if err != nil {
		log.Fatalf("failed to create output file: %v", err)
	}
	defer f.Close()

	t := template.Must(template.New("iso4217").Parse(tmpl))
	if err := t.Execute(f, Data{N2C: n2c, C2N: c2n}); err != nil {
		log.Fatalf("failed to execute template: %v", err)
	}
}
