package main

import "testing"

func TestQueryURL(t *testing.T) {
	s1 := queryUrl("rotwein", "0")
	if s1 != "https://www.chefkoch.de/rs/s0/rotwein/Rezepte.html#more2" {
		t.Error("Expected s1 to be , got: ", s1)
	}
}
