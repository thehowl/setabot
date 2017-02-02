// Package stops provide useful structs of stops in the cities.
package stops

// Stop represents a bus stop.
type Stop struct {
	Name       string
	Identifier string
}

// CityStops is a map holding the various stops for each city.
var CityStops = map[string][]Stop{
	"mo": Mo,
	"re": Re,
	"pc": Pc,
}
