package scrape

// This contains the errors used throughout the scrape package

import "errors"

var (
	//ErrInvalidID is returned when a scrape id is found to be invalid
	ErrInvalidID = errors.New("Invalid ID, please ignore in the future")
	//ErrUniversityPanicking is returned if the ical server at the university returns a 5xx status
	ErrUniversityPanicking = errors.New("The uni server seems to be panicking, slow down and let them catch their breath")
	//ErrConfig is returned if the scraper configuration is invalid
	ErrConfig = errors.New("Scraper configuration invalid")
)
