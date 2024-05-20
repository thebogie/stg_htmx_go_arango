package types

import "sort"

type Game struct {
	Id   string `json:"_id"`
	Name string `json:"name"`
}

type FindGameAPI struct {
	FindGame []Game `json:"FindGame"`
}

func (f *FindGameAPI) SortByName() {
	sort.Slice(f.FindGame, func(i, j int) bool {
		return f.FindGame[i].Name < f.FindGame[j].Name
	})
}
