package types

import "sort"

type Game struct {
	Id            string `json:"_id"`
	Name          string `json:"name"`
	YearPublished int    `json:"year_published"`
	BGGId         int    `json:"bgg_id"`
}

type Games struct {
	List []Game
}

func (f *Games) SortByName() {
	sort.Slice(f.List, func(i, j int) bool {
		return f.List[i].Name < f.List[j].Name
	})
}

type FindGameAPI struct {
	FindGame []Game `json:"FindGame"`
}

func (f *FindGameAPI) SortByName() {
	sort.Slice(f.FindGame, func(i, j int) bool {
		return f.FindGame[i].Name < f.FindGame[j].Name
	})
}
