package scraper_test

import (
	"testing"

	"github.com/thehowl/setabusbot/providers/scraper"
)

func TestGetArrivals(t *testing.T) {
	s := new(scraper.Scraper)
	arrs, err := s.GetArrivals("mo", "MO10", "MODENA  AUTOSTAZIONE")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(arrs)
}
