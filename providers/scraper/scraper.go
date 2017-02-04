// Package scraper implements setabusbot's services through web scraping.
package scraper

import (
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/thehowl/setabusbot/services"
)

// Scraper is a provider of the services of setabusbot that fetches data through
// web scraping.
type Scraper struct{}

const timeFormat = "15:04"

// GetArrivals fetches the arrivals at a stop.
func (s *Scraper) GetArrivals(city, stopID, stopName string) ([]services.Arrival, error) {
	// create POST values and make the request
	vals := make(url.Values)
	vals.Add("risultato", "palina")
	vals.Add("nome_fermata", stopName)
	vals.Add("qm_palina", stopID)
	vals.Add("x", "13")
	vals.Add("y", "12")
	resp, err := http.PostForm("http://www.setaweb.it/mo/quantomanca", vals)
	if err != nil {
		return nil, err
	}

	// create goquery and start looking for what we need
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return nil, err
	}
	rows := doc.Find("table.qm_table_risultati").First().Find("tbody tr")

	// create arrivals slice which will then be returned
	arrivals := make([]services.Arrival, 0, rows.Length())

	// loop over the rows to get all the arrivals.
	rows.Each(func(_ int, s *goquery.Selection) {
		var a services.Arrival
		s.Find("td").Each(func(idx int, s *goquery.Selection) {
			// based on the position of the td, we can understand what
			// information it gives us.
			switch idx {
			case 0:
				a.Line = s.Text()
			case 1:
				a.Destination = s.Text()
			case 2:
				a.Urban = !strings.Contains(s.Children().First().AttrOr("src", ""), "icona_linea_extraurbana")
			case 3:
				a.ToArrival = strings.TrimSuffix(s.Text(), "'")
			case 4:
				a.TimetableTime, _ = time.Parse(timeFormat, s.Text())
			case 5:
				a.RealTime, _ = time.Parse(timeFormat, s.Text())
			}
		})
		if a.Destination != "" {
			arrivals = append(arrivals, a)
		}
	})

	return arrivals, nil
}
