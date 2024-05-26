package view

import (
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
	suggestionResult, err := app.SuggestBook(pref)
	if err != nil {
		SparkResult([]app.BookSpark{}, nil, nil, err).Render(r.Context(), w)
	} else {
		for title, analysis := range suggestionResult.BookAnalysis {
			pieChartLabels = append(pieChartLabels, title)
			pieChartData = append(pieChartData, analysis.Count)
		}
		SparkResult(suggestionResult.Books, pieChartLabels, pieChartData, nil).Render(r.Context(), w)
	}
}
