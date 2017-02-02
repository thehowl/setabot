// +build ignore

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var cities = []string{
	"mo",
	"pc",
	"re",
}

func main() {
	err := execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func execute() error {
	for _, c := range cities {
		err := scrapeCity(c)
		if err != nil {
			return err
		}
	}
	return nil
}

const header = `package stops

// THIS FILE HAS BEEN AUTOMATICALLY GENERATED!
// Check out scripts/stop_scraper.go

var %s = []Stop{
`

func scrapeCity(c string) error {
	f, err := os.Create("stops/" + c + ".go")
	if err != nil {
		return err
	}
	defer f.Close()
	// write header
	f.Write([]byte(fmt.Sprintf(header, strings.Title(c))))

	// scrape seta website for this city
	doc, err := goquery.NewDocument("http://www.setaweb.it/" + c + "/quantomanca")
	if err != nil {
		return err
	}

	// get the stuff and add it to the file.
	doc.Find("#qm_palina").Children().Each(func(_ int, s *goquery.Selection) {
		if s.Text() == "" {
			return
		}
		// write to file new bus top
		f.Write([]byte(fmt.Sprintf("\t{%q, %q},\n", strings.TrimSpace(s.Text()), s.AttrOr("value", ""))))
	})

	f.Write([]byte("}\n"))

	return nil
}
