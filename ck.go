package main

import (
	"encoding/json"
	"google.golang.org/appengine"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Recipe struct {
	Title      string `json:"title"`
	Subtitle   string `json:"subtitle"`
	Url        string `json:"url"`
	Thumbnail  string `json:"thumbnail"`
	Rating     string `json:"rating"`
	Difficulty string `json:"difficulty"`
	Preptime   string `json:"preptime"`
}

type RecipeDetail struct {
	Title       string              `json:"title"`
	Rating      string              `json:"rating"`
	Difficulty  string              `json:"difficulty"`
	Preptime    string              `json:"preptime"`
	Cookingtime string              `json:"cookingtime"`
	Thumbnail   string              `json:"thumbnail"`
	Ingredients []*RecipeIngredient `json:"ingredients"`
	Method      string              `json:"method"`
}

type RecipeIngredient struct {
	Amount     string `json:"amount"`
	Ingredient string `json:"ingredient"`
}

type RecipeDetailDocument struct {
	doc *goquery.Document
}

func (rdd *RecipeDetailDocument) newRecipeDetail() *RecipeDetail {
	title := rdd.title()
	ingredients := rdd.ingredients()
	method := rdd.method()
	rating := rdd.rating()
	prepinfo := rdd.prepinfo()
	difficulty := rdd.difficulty(prepinfo)
	preptime := rdd.preptime(prepinfo)
	cookingtime := rdd.cookingtime(prepinfo)
	thumbnail := rdd.thumbnail()
	return &RecipeDetail{title, rating, difficulty,
		preptime, cookingtime, thumbnail, ingredients, method}
}

const CKPrefix = "https://www.chefkoch.de"

func queryUrl(searchterm string, page string) string {
	return "https://www.chefkoch.de/rs/s" + page + "/" + searchterm + "/Rezepte.html#more2"
}

type RecipesSelection struct {
	sel *goquery.Selection
}

func (rs *RecipesSelection) title() string {
	return rs.sel.Find(".search-list-item-title").Text()
}

func (rs *RecipesSelection) subtitle() string {
	subtitle := rs.sel.Find(".search-list-item-subtitle").Text()
	trimmed := strings.Trim(subtitle, " \n")
	return strings.Replace(trimmed, "\n", " ", -1)
}

func (rs *RecipesSelection) url() string {
	return CKPrefix + rs.sel.Find(".search-list-item > a").AttrOr("href", "")
}

func (rs *RecipesSelection) thumbnail() string {
	thumb := rs.sel.Find("picture > img").AttrOr("srcset", "")
	if strings.HasPrefix(thumb, "data:image") {
		thumb = rs.sel.Find("picture > img").AttrOr("data-srcset", "")
	}
	return thumb
}

func (rs *RecipesSelection) rating() string {
	rating := rs.sel.Find(".search-list-item-uservotes-stars").AttrOr("title", "")
	digregex := regexp.MustCompile(`\d\.\d*`)
	return digregex.FindString(rating)

}

func (rdd *RecipeDetailDocument) title() string {
	return rdd.doc.Find(".page-title").Text()
}

func (rdd *RecipeDetailDocument) rating() string {
	rat := rdd.doc.Find(".rating__average-rating").Text()
	rat = strings.Replace(rat, "Ø", "", 1)
	return strings.Replace(rat, ",", ".", 1)
}

func (rdd *RecipeDetailDocument) difficulty(pi map[string]string) string {
	return pi["Schwierigkeitsgrad"]
}

func (rdd *RecipeDetailDocument) preptime(pi map[string]string) string {
	return pi["Arbeitszeit"]
}

func (rdd *RecipeDetailDocument) cookingtime(pi map[string]string) string {
	return pi["Kochzeit"]
}

func (rdd *RecipeDetailDocument) prepinfo() map[string]string {
	prep := rdd.doc.Find("#preparation-info").Text()
	prep = strings.Replace(prep, "\n", "", -1)
	prep = strings.Replace(prep, "Koch-/Backzeit", "Kochzeit", 1)
	prep = strings.Replace(prep, "keine Angabe", "NA", 1)
	sections := strings.Split(prep, "/")
	result := make(map[string]string)
	for _, i := range sections {
		trimmed := strings.Trim(i, " \n\t")
		split := strings.Split(trimmed, ":")
		result[split[0]] = strings.Trim(split[1], " \n")
	}
	return result
}

func (rdd *RecipeDetailDocument) thumbnail() string {
	return rdd.doc.Find(".slideshow-image").AttrOr("src", "")
}

func (rdd *RecipeDetailDocument) ingredients() []*RecipeIngredient {
	ingtable := rdd.doc.Find(".incredients>tbody>tr")
	var ingredients []*RecipeIngredient
	ingtable.Each(func(i int, s *goquery.Selection) {
		amount := strings.Trim(s.Find(".amount").Text(), " \n")
		ing := strings.Trim(s.Find("td:nth-child(2)").Text(), " \n")
		ingredients = append(ingredients, &RecipeIngredient{amount, ing})
	})
	return ingredients
}

func (rdd *RecipeDetailDocument) method() string {
	text := rdd.doc.Find("#rezept-zubereitung").Text()
	return strings.Trim(text, " \n")
}

func (rs *RecipesSelection) difficulty() string {
	return rs.sel.Find(".search-list-item-difficulty").Text()
}

func (rs *RecipesSelection) preptime() string {
	return rs.sel.Find(".search-list-item-preptime").Text()
}

func NewRecipe(sel *goquery.Selection) *Recipe {
	rs := &RecipesSelection{sel}
	return &Recipe{rs.title(), rs.subtitle(),
		rs.url(), rs.thumbnail(), rs.rating(), rs.difficulty(),
		rs.preptime()}
}

func allRecipes(doc *goquery.Document) []*Recipe {
	var results []*Recipe
	doc.Find(".search-list-item").Each(func(i int, s *goquery.Selection) {
		results = append(results, NewRecipe(s))
	})
	return results
}

func recipesToJson(recipes []*Recipe) ([]byte, error) {
	return json.Marshal(recipes)
}

func recipeDetailToJson(rd *RecipeDetail) ([]byte, error) {
	return json.Marshal(rd)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	query := r.FormValue("query")
	page := r.FormValue("page")
	url := queryUrl(query, page)
	res, err := http.Get(url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	json, err := recipesToJson(allRecipes(doc))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write(json)
}

func detailHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	recurl := r.FormValue("recipeurl")
	res, err := http.Get(recurl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	rddoc := RecipeDetailDocument{doc}
	rdd := rddoc.newRecipeDetail()
	json, err := recipeDetailToJson(rdd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write(json)
}
func main() {
	http.HandleFunc("/search", searchHandler)
	http.HandleFunc("/recipedetail", detailHandler)
	appengine.Main()
}
