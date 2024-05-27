package app

import (
	"errors"
	"fmt"
	"math"
	"sort"
	"strings"
)

type Preferences struct {
	Language       string
	Genre          string
	TargetAudience string
	Subject        string
}

type BookSpark struct {
	Title                     string
	Author                    string
	PreferenceMatchPercentage int
}

type BookAnalysis struct {
	Book                      Book
	Count                     int
	PreferenceMatchPercentage int
}

type SuggestionResult struct {
	Books        []BookSpark
	BookAnalysis map[string]*BookAnalysis
}

func SuggestBook(pref Preferences) (SuggestionResult, error) {
	var books []BookSpark
	openlibClient := NewOpenLibClient()
	limit := 10
	qResult, err := openlibClient.OpenLib.Query(pref.Language, pref.Genre, pref.TargetAudience, pref.Subject, limit)
	if err != nil {
		fmt.Printf("Error during query: ", err)
	}
	initialBooks, err := openlibClient.OpenLib.GetBooks(qResult.Docs)
	if err != nil {
		fmt.Printf("Error during get books: ", err)
	}

	if len(initialBooks) == 0 {
		return SuggestionResult{}, errors.New("no books found or service unavailable")
	}

	bookAnalysis := make(map[string]*BookAnalysis)
	for _, book := range initialBooks {
		analysis := &BookAnalysis{Book: book, Count: 0}
		prefMatchCount := 0
		for _, subject := range book.Subjects {
			if strings.Contains(strings.ToLower(subject), strings.ToLower(getLanguageName(pref.Language))) ||
				strings.Contains(strings.ToLower(subject), strings.ToLower(pref.Genre)) ||
				strings.Contains(strings.ToLower(subject), strings.ToLower(pref.Subject)) ||
				strings.Contains(strings.ToLower(subject), strings.ToLower(pref.TargetAudience)) {
				prefMatchCount++
			}
		}
		if len(book.Subjects) > 0 {
			percentage := (float64(prefMatchCount) / float64(len(book.Subjects))) * 100
			analysis.PreferenceMatchPercentage = int(math.Ceil(percentage))
		}
		bookAnalysis[book.Title] = analysis
	}

	for i := 0; i < 10; i++ {
		newBooks, err := openlibClient.OpenLib.GetBooks(qResult.Docs)
		if err != nil {
			fmt.Println("Error fetching data:", err)
			continue
		}
		if len(newBooks) == 0 {
			fmt.Println("No books found or error fetching data: ", err)
			continue
		}

		for _, ba := range bookAnalysis {
			if ba.Book.Title == newBooks[0].Title {
				ba.Count++
			}
		}
	}

	bookAnalysisSlice := make([]*BookAnalysis, 0, len(bookAnalysis))
	for _, ba := range bookAnalysis {
		bookAnalysisSlice = append(bookAnalysisSlice, ba)
	}

	sort.Slice(bookAnalysisSlice, func(i, j int) bool {
		if bookAnalysisSlice[i].PreferenceMatchPercentage == bookAnalysisSlice[j].PreferenceMatchPercentage {
			return bookAnalysisSlice[i].Count > bookAnalysisSlice[j].Count
		}
		return bookAnalysisSlice[i].PreferenceMatchPercentage > bookAnalysisSlice[j].PreferenceMatchPercentage
	})

	for _, analysis := range bookAnalysisSlice {
		book := BookSpark{
			Title:                     analysis.Book.Title,
			Author:                    analysis.Book.Author,
			PreferenceMatchPercentage: analysis.PreferenceMatchPercentage,
		}
		books = append(books, book)
	}
	return SuggestionResult{Books: books, BookAnalysis: bookAnalysis}, err
}

func getLanguageName(code string) string {
	if code == "eng" {
		return "English"
	}
	if code == "spa" {
		return "Spanish"
	}
	if code == "fre" {
		return "French"
	}
	if code == "deu" {
		return "German"
	}
	if code == "zho" {
		return "Chinese"
	}
	if code == "jpn" {
		return "Japanese"
	}
	return "English"
}
