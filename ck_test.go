package main

import "testing"

var queryurls = []struct {
	st   string
	page string
	want string
}{
	{
		"rotwein",
		"0",
		"https://www.chefkoch.de/rs/s0/rotwein/Rezepte.html#more2",
	},
	{
		"rotwein",
		"60",
		"https://www.chefkoch.de/rs/s60/rotwein/Rezepte.html#more2",
	},
	{
		"bohnen",
		"120",
		"https://www.chefkoch.de/rs/s120/bohnen/Rezepte.html#more2",
	},
}

func TestQueryURL(t *testing.T) {
	for _, i := range queryurls {
		qu := queryUrl(i.st, i.page)
		if qu != i.want {
			t.Errorf("Expected qu to be %s, got: %s", i.want, qu)
		}
	}
}
