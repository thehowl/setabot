package services

import "time"

// ArrivalsService is a service which can fetch the bus arrival times of a
// certain stop.
type ArrivalsService interface {
	GetArrivals(city, stopID, stopName string) ([]Arrival, error)
}

// Arrival represents an arrival at our stop.
type Arrival struct {
	Line          string
	Destination   string
	Urban         bool
	ToArrival     string    // Minutes missing before the bus arrives
	TimetableTime time.Time // The time of arrival, as written on the timetable
	RealTime      time.Time // The real time of arrival, including delays
}
