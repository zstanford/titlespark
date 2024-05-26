package view

import (
	"fmt"
	"net/http"
	"titlespark-web/internal/app"

	"github.com/go-chi/chi/v5"
)

func Handlers() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/", index)
	r.Get("/spark", sparkForm)
	r.Post("/spark", sparkSuggestion)
	return r
}

func index(w http.ResponseWriter, r *http.Request) {
	Spark().Render(r.Context(), w)
}

func sparkForm(w http.ResponseWriter, r *http.Request) {
	SparkForm().Render(r.Context(), w)
}

func sparkSuggestion(w http.ResponseWriter, r *http.Request) {
	var (
		pieChartLabels []string
		pieChartData   []int
	)
	r.ParseForm()
	pref := app.Preferences{
		Language:       r.PostForm.Get("language"),
		Genre:          r.PostForm.Get("genre"),
		TargetAudience: r.PostForm.Get("target-audience"),
		Subject:        r.PostForm.Get("subject"),
	}
	fmt.Printf("Language: %s\nGenre: %s\nTarget Audience: %s\nSubject: %s\n",
		pref.Language, pref.Genre, pref.TargetAudience, pref.Subject)
	books, err := app.SuggestBook(pref)
	if err != nil {
		SparkResult([]app.BookSpark{}, nil, nil, err).Render(r.Context(), w)
	} else {
		for index, book := range books {
			pieChartLabels = append(pieChartLabels, book.Title)
			pieChartData = append(pieChartData, index+1)
		}
		SparkResult(books, pieChartLabels, pieChartData, nil).Render(r.Context(), w)
	}
}
