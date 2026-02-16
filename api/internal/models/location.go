package models

type Location struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	ParkCode  string  `json:"park_code"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	State     string  `json:"state"`
	Region    string  `json:"region"`
}
