package main

import (
	"bytes"
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"reflect"
	"testing"
)

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

var bohnenrecipes = []struct {
	title      string
	subtitle   string
	url        string
	thumbnail  string
	rating     string
	difficulty string
	preptime   string
}{
	{
		"Grüne Bohnen im Speckmantel",
		"Bohnen waschen und die Spitzen abschneiden. Bohnenkraut, Knoblauch, zerdrückte Pfefferkörner und Salz mit Öl kurz anrösten. 2 Lite...",
		"https://www.chefkoch.de/rezepte/563451154612271/Gruene-Bohnen-im-Speckmantel.html",
		"https://static.chefkoch-cdn.de/rs/bilder/56345/gruene-bohnen-im-speckmantel-1124631-150x150.jpg",
		"4.49",
		"simpel",
		"30 min.",
	},
	{
		"Grüne Bohnen",
		"Variante 1: Die Bohnen putzen. Zwiebeln und Knoblauch klein schneiden und in etwas Butter oder Margarine anbraten. Die Bohnen dazu...",
		"https://www.chefkoch.de/rezepte/3166211471333987/Gruene-Bohnen.html",
		"https://static.chefkoch-cdn.de/rs/bilder/316621/gruene-bohnen-938192-150x150.jpg",
		"4.36",
		"simpel",
		"10 min.",
	},
	{
		"Schupfnudel - Bohnen - Pfanne",
		"Pfannengericht mit Bohnen, Schinken, Schupfnudeln und Crème fraiche",
		"https://www.chefkoch.de/rezepte/1171381223217983/Schupfnudel-Bohnen-Pfanne.html",
		"https://static.chefkoch-cdn.de/rs/bilder/117138/schupfnudel-bohnen-pfanne-1156413-150x150.jpg",
		"4.37",
		"normal",
		"30 min.",
	},
	{
		"Grüne Bohnen mit Speck",
		"Speckbohnen",
		"https://www.chefkoch.de/rezepte/2406611380140966/Gruene-Bohnen-mit-Speck.html",
		"https://static.chefkoch-cdn.de/rs/bilder/240661/gruene-bohnen-mit-speck-1135575-150x150.jpg",
		"4.67",
		"normal",
		"25 min.",
	},
	{
		"Grüne Bohnen mit Kasseler, geschmort",
		"Zwiebeln in Butter anbraten. Kasselerwürfel dazu geben und ebenfalls anbraten. Die Grünen Bohnen (Menge je nach Geschmack) in Stüc...",
		"https://www.chefkoch.de/rezepte/103621042299597/Gruene-Bohnen-mit-Kasseler-geschmort.html",
		"https://static.chefkoch-cdn.de/rs/bilder/10362/gruene-bohnen-mit-kasseler-geschmort-1135594-150x150.jpg",
		"4.5",
		"normal",
		"30 min.",
	},
}

var sahneRecipes = []struct {
	title      string
	subtitle   string
	url        string
	thumbnail  string
	rating     string
	difficulty string
	preptime   string
}{
	{},
}

func TestNewRecipe(t *testing.T) {
	file, err := ioutil.ReadFile("testhtml/bohnen.html")
	if err != nil {
		t.Error(err)
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(file))
	if err != nil {
		panic(err)
	}
	var results []*Recipe
	doc.Find(".search-list-item").Each(func(i int, sel *goquery.Selection) {
		results = append(results, NewRecipe(sel))
	})
	if len(results) != 30 {
		t.Error("Expected length of results to be 30, got: ", len(results))
	}
	for ind, i := range bohnenrecipes {
		if results[ind].Title != i.title {
			t.Errorf("Expected title to be %q, got: %q", i.title, results[ind].Title)
		}
		if results[ind].Subtitle != i.subtitle {
			t.Errorf("Expected subtitle to be %q, got: %q", i.subtitle,
				results[ind].Subtitle)
		}
		if results[ind].Url != i.url {
			t.Errorf("Expected url to be %q, got: %q", i.url,
				results[ind].Url)
		}
		if results[ind].Thumbnail != i.thumbnail {
			t.Errorf("Expected thumbnail to be %q, got: %q", i.thumbnail,
				results[ind].Thumbnail)
		}
		if results[ind].Rating != i.rating {
			t.Errorf("Expected rating to be %q, got %q", i.rating,
				results[ind].Rating)
		}
		if results[ind].Difficulty != i.difficulty {
			t.Errorf("Expected difficulty to be %q, got %q", i.difficulty,
				results[ind].Difficulty)
		}
		//if results[ind].Preptime != i.preptime {
		//	t.Errorf("Expected preptime to be %q, got %q", i.preptime,
		//		results[ind].Preptime)
		//}

	}
}

func TestRecipesToJSON(t *testing.T) {
	file, err := ioutil.ReadFile("testhtml/bohnen.html")
	if err != nil {
		panic(err)
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(file))
	if err != nil {
		panic(err)
	}
	recipes := allRecipes(doc)
	jsonrecipes, err := recipesToJson(recipes)
	if err != nil {
		t.Error("Expected error to be nil, got: ", err)
	}
	var unmarshalled []*Recipe
	json.Unmarshal(jsonrecipes, &unmarshalled)
	if unmarshalled[0].Title != bohnenrecipes[0].title {
		t.Errorf("Expected Title to be %q, got %q",
			bohnenrecipes[0].title, unmarshalled[0].Title)
	}
}

var grueneImSpeckmantel = struct {
	title       string
	thumbnail   string
	ingredients []*RecipeIngredient
	method      string
	rating      string
	difficulty  string
	cookingtime string
	preptime    string
}{

	"Grüne Bohnen im Speckmantel",
	"https://static.chefkoch-cdn.de/ck.de/rezepte/56/56345/1124631-420x280-fix-gruene-bohnen-im-speckmantel.jpg",
	[]*RecipeIngredient{
		{"800\u00a0g", "Bohnen, frische"},
	},
	"Bohnen waschen und die Spitzen abschneiden.\nBohnenkraut, Knoblauch, zerdrückte Pfefferkörner und Salz mit Öl kurz anrösten. 2 Liter Wasser zugießen, 10 Min. kochen, durchsieben. Diese Brühe aufkochen und die Bohnen in 3 Portionen nacheinander sprudelnd garen. Schnell in kaltem Wasser abkühlen, in einem Tuch abtrocknen.\n\nBohnen in Bacon einwickeln. Butter in einer feuerfesten Form erhitzen, die Bohnen reingeben (mit der Specknaht nach unten) und zugedeckt im Ofen bei 180 °C - 200 °C erhitzen (ca. 5 Minuten), dabei einmal wenden.",
	"4.49",
	"simpel",
	"ca. 15 Min.",
	"ca. 30 Min.",
}

var schupfnudel = struct {
	title       string
	thumbnail   string
	ingredients []*RecipeIngredient
	method      string
	rating      string
	difficulty  string
	cookingtime string
	preptime    string
}{
	"Schupfnudel - Bohnen - Pfanne",
	"https://static.chefkoch-cdn.de/ck.de/rezepte/117/117138/1156413-420x280-fix-schupfnudel-bohnen-pfanne.jpg",
	[]*RecipeIngredient{
		{"500\u00a0g", "Schupfnudeln (Kühlregal)"},
		{"200\u00a0g", "Schinken, gekochter"},
		{"250\u00a0g", "Bohnen (Prinzessbohnen, TK)"},
		{"1/8\u00a0Liter", "Fleischbrühe"},
		{"1\u00a0Becher", "Crème fraîche"},
		{"4\u00a0Scheibe/n", "Käse (Toast-Käse, z.B. Scheibletten)"},
		{"\u00a0n. B.", "Salz und Pfeffer"},
		{"\u00a0", "Olivenöl"},
	},
	"Die Prinzessböhnchen für ca. 5 Min. in kochendem Wasser garen. \n\nDen Kochschinken würfeln und mit etwas Olivenöl in der Pfanne anbraten. Die Schupfnudeln hinzugeben und 5-8 Min. zusammen mit dem Schinken braten, bis die Schupfnudeln eine goldgelbe Farbe annehmen. Die Prinzessbohnen hinzu geben. Nun 1/8 l Fleischbrühe zugießen und mit Crème fraiche nach Belieben andicken. Nach Geschmack würzen. Als Abschluss die Käsescheiben oben auflegen, bis diese verlaufen. Sofort servieren.",
	"4.37",
	"normal",
	"",
	"ca. 30 Min.",
}

var speckbohnen = struct {
	title       string
	thumbnail   string
	ingredients []*RecipeIngredient
	method      string
	rating      string
	difficulty  string
	cookingtime string
	preptime    string
}{
	"Grüne Bohnen mit Speck",
	"https://static.chefkoch-cdn.de/ck.de/rezepte/240/240661/1135575-420x280-fix-gruene-bohnen-mit-speck.jpg",
	[]*RecipeIngredient{
		{"500\u00a0g", "Bohnen, grüne, frisch oder TK"},
	},
	"Grüne Bohnen putzen, ca. 5 Min. in Salzwasser ankochen (bei TK nach Anleitung kochen). Speck würfeln und im Butter-Öl Gemisch kross anbraten. Bohnen, Speck und Bohnenkraut zusammen in einen Topf geben, pfeffern und 10-20 Min. bei kleiner Hitze ziehen lassen, gelegentlich umrühren. Wem es zu kräftig (salzig) ist, einfach weniger Speck nehmen.",
	"4.67",
	"normal",
	"ca. 20 Min.",
	"ca. 25 Min.",
}

func TestNewRecipeDetail(t *testing.T) {
	detailfile, err := ioutil.ReadFile("testhtml/gruene_bohnen_im_speckmantel.html")
	if err != nil {
		panic(err)
	}
	detaildoc, err := goquery.NewDocumentFromReader(bytes.NewReader(detailfile))
	if err != nil {
		panic(err)
	}
	rdd := &RecipeDetailDocument{detaildoc}
	grbohndetail := rdd.newRecipeDetail()
	if grbohndetail.Title != grueneImSpeckmantel.title {
		t.Errorf("Expected title to be %q, got %q", grueneImSpeckmantel.title,
			grbohndetail.Title)
	}
	if grbohndetail.Ingredients[0].Amount != grueneImSpeckmantel.ingredients[0].Amount {
		t.Errorf("Expected amount to be '800 g', got: %q",
			grbohndetail.Ingredients[0].Amount)
	}
	if grbohndetail.Ingredients[0].Ingredient != grueneImSpeckmantel.ingredients[0].Ingredient {
		t.Errorf("Expected ingredient to be 'Bohnen, frische', got: %q",
			grbohndetail.Ingredients[0].Ingredient)
	}
	if grbohndetail.Method != grueneImSpeckmantel.method {
		t.Errorf("Expected method to be %q, got: %q",
			grueneImSpeckmantel.method, grbohndetail.Method)
	}
	if grbohndetail.Difficulty != grueneImSpeckmantel.difficulty {
		t.Errorf("Expected difficulty to be %q, got: %q",
			grueneImSpeckmantel.difficulty, grbohndetail.Difficulty)
	}
	if grbohndetail.Thumbnail != grueneImSpeckmantel.thumbnail {
		t.Errorf("Expected thumbnail to be %q, got: %q",
			grueneImSpeckmantel.thumbnail, grbohndetail.Thumbnail)
	}
	if grbohndetail.Rating != grueneImSpeckmantel.rating {
		t.Errorf("Expected rating to be %q, got: %q",
			grueneImSpeckmantel.rating, grbohndetail.Rating)
	}
	if grbohndetail.Preptime != grueneImSpeckmantel.preptime {
		t.Errorf("Expected preptime to be %q, got: %q",
			grueneImSpeckmantel.preptime, grbohndetail.Preptime)
	}
	if grbohndetail.Cookingtime != grueneImSpeckmantel.cookingtime {
		t.Errorf("Expected cookingtime to be %q, got: %q",
			grueneImSpeckmantel.cookingtime, grbohndetail.Cookingtime)
	}
	detailfile, err = ioutil.ReadFile("testhtml/schupfnudel.html")
	if err != nil {
		panic(err)
	}
	detaildoc, err = goquery.NewDocumentFromReader(bytes.NewReader(detailfile))
	if err != nil {
		panic(err)
	}
	rdd = &RecipeDetailDocument{detaildoc}
	schupfdetail := rdd.newRecipeDetail()
	if schupfdetail.Title != schupfnudel.title {
		t.Errorf("Expected title to be %q, got %q", schupfnudel.title,
			schupfdetail.Title)
	}
	for ind, i := range schupfdetail.Ingredients {
		if i.Amount != schupfnudel.ingredients[ind].Amount {
			t.Errorf("Expected amount to be %q, got %q",
				schupfnudel.ingredients[ind].Amount, i.Amount)
		}
	}

	for ind, i := range schupfdetail.Ingredients {
		if i.Ingredient != schupfnudel.ingredients[ind].Ingredient {
			t.Errorf("Expected ingredient to be %q, got %q",
				schupfnudel.ingredients[ind].Ingredient, i.Ingredient)
		}
	}
	if schupfdetail.Method != schupfnudel.method {
		t.Errorf("Expected method to be %q, got %q",
			schupfnudel.method, schupfdetail.Method)
	}
	if schupfdetail.Thumbnail != schupfnudel.thumbnail {
		t.Errorf("Expected thumbnail to be %q, got %q",
			schupfnudel.thumbnail, schupfdetail.Thumbnail)
	}
	if schupfdetail.Difficulty != schupfnudel.difficulty {
		t.Errorf("Expected difficulty to be %q, got %q",
			schupfnudel.difficulty, schupfdetail.Difficulty)
	}
	if schupfdetail.Rating != schupfnudel.rating {
		t.Errorf("Expected rating to be %q, got %q",
			schupfnudel.rating, schupfdetail.Rating)
	}
	if schupfdetail.Preptime != schupfnudel.preptime {
		t.Errorf("Expected preptime to be %q, got %q",
			schupfnudel.preptime, schupfdetail.Preptime)
	}
	if schupfdetail.Cookingtime != schupfnudel.cookingtime {
		t.Errorf("Expected cooking time to be %q, got %q",
			schupfnudel.cookingtime, schupfdetail.Cookingtime)
	}
	detailfile, err = ioutil.ReadFile("testhtml/gruene_bohnen_mit_speck.html")
	if err != nil {
		panic(err)
	}
	detaildoc, err = goquery.NewDocumentFromReader(bytes.NewReader(detailfile))
	if err != nil {
		panic(err)
	}
	rdd = &RecipeDetailDocument{detaildoc}
	spdetail := rdd.newRecipeDetail()
	if spdetail.Title != speckbohnen.title {
		t.Errorf("Expected title to be %q, got %q",
			speckbohnen.title, spdetail.Title)
	}
	if spdetail.Ingredients[0].Amount != speckbohnen.ingredients[0].Amount {
		t.Errorf("Expected amount to be %q, got %q",
			speckbohnen.ingredients[0].Amount, spdetail.Ingredients[0].Amount)
	}
	if spdetail.Ingredients[0].Ingredient != speckbohnen.ingredients[0].Ingredient {
		t.Errorf("Expected ingredient to be %q, got %q",
			speckbohnen.ingredients[0].Ingredient, spdetail.Ingredients[0].Ingredient)
	}
	if spdetail.Method != speckbohnen.method {
		t.Errorf("Epected method to be %q, got \n%q",
			speckbohnen.method, spdetail.Method)
	}
	if spdetail.Difficulty != speckbohnen.difficulty {
		t.Errorf("Expected difficulty to be %q, got %q",
			speckbohnen.difficulty, spdetail.Difficulty)
	}
	if spdetail.Thumbnail != speckbohnen.thumbnail {
		t.Errorf("Expected thumbnail to be %q, got %q",
			speckbohnen.thumbnail, spdetail.Thumbnail)
	}
	if spdetail.Rating != speckbohnen.rating {
		t.Errorf("Expected rating to be %q, got %q",
			speckbohnen.rating, spdetail.Rating)
	}
	if spdetail.Preptime != speckbohnen.preptime {
		t.Errorf("Expected preptime to be %q, got %q",
			speckbohnen.preptime, spdetail.Preptime)
	}
	if spdetail.Cookingtime != speckbohnen.cookingtime {
		t.Errorf("Expected cookingtime to be %q, got %q",
			speckbohnen.cookingtime, spdetail.Cookingtime)
	}
}

func TestNewRecipeDetailToJson(t *testing.T) {
	file, err := ioutil.ReadFile("testhtml/schupfnudel.html")
	if err != nil {
		panic(err)
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(file))
	if err != nil {
		panic(err)
	}
	rdd := &RecipeDetailDocument{doc}
	rd := rdd.newRecipeDetail()
	marschalled, err := recipeDetailToJson(rd)
	if err != nil {
		panic(err)
	}
	var nrd *RecipeDetail
	err = json.Unmarshal(marschalled, &nrd)
	if err != nil {
		panic(err)
	}
	if nrd.Cookingtime != rd.Cookingtime {
		t.Errorf("Expected cookingtime of unmarschalled to be equal to original, got: %q",
			nrd.Cookingtime)
	}
	if nrd.Method != rd.Method {
		t.Error("Expected unmarschalled method to be equal to original, got: ",
			nrd.Method)
	}
	if !reflect.DeepEqual(nrd, rd) {
		t.Error("ecpected both interfaces to be equal.")
	}

}
