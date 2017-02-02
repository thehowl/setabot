package scraper_test

import (
	"testing"

	"github.com/thehowl/setabusbot/providers/scraper"
)

func TestGetArrivals(t *testing.T) {
	s := new(scraper.Scraper)
	arrs, err := s.GetArrivals("mo", "MO3745", "MARZAGLIA VECCHIA")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(arrs)
}
