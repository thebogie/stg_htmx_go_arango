package types

import "sort"

type Venue struct {
	Id      string  `json:"_id"`
	Address string  `json:"address"`
	PlaceID string  `json:"place_id"`
	Lat     float64 `json:"lat"`
	Lng     float64 `json:"lng"`
}

type Venues struct {
	List []Venue
}

func (f *Venues) SortByAddress() {
	sort.Slice(f.List, func(i, j int) bool {
		return f.List[i].Address < f.List[j].Address
	})
}

type FindVenueAPI struct {
	FindVenue []Venue `json:"FindVenue"`
}
