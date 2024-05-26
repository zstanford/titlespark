package app

import (
	"errors"
	"fmt"
)

type Preferences struct {
	Language       string
	Genre          string
	TargetAudience string
	Subject        string
}

type BookSpark struct {
	Title  string
	Author string
}

type BookAnalysis struct {
	Book  Book
	Count int
}

type SuggestionResult struct {
	Books        []BookSpark
	BookAnalysis map[string]*BookAnalysis
}

func SuggestBook(pref Preferences) (SuggestionResult, error) {
	var books []BookSpark
	openlibClient := NewOpenLibClient()
	qResult, err := openlibClient.OpenLib.Query(pref.Language, pref.Genre, pref.TargetAudience, pref.Subject)
	if err != nil {
		fmt.Printf("Error during query: ", err)
	}
	initialBooks, err := openlibClient.OpenLib.GetBooks(qResult.Docs)
	if err != nil {
		fmt.Printf("Error during get books: ", err)
	}

	if len(initialBooks) == 0 {
		return SuggestionResult{}, errors.New("no books found")
	}
	fmt.Printf("Book Result length: %v\n", len(initialBooks))

	bookAnalysis := make(map[string]*BookAnalysis)
	for _, book := range initialBooks {
		bookAnalysis[book.Title] = &BookAnalysis{Book: book, Count: 0}
	}

	for i := 0; i < 10; i++ {
		newBooks, err := openlibClient.OpenLib.GetBooks(qResult.Docs)
		if err != nil {
			fmt.Println("Error fetching data:", err)
			continue
		}
		for _, ba := range bookAnalysis {
			if ba.Book.Title == newBooks[0].Title {
				ba.Count++
			}
		}
	}
	for _, analysis := range bookAnalysis {
		book := BookSpark{
			Title:  analysis.Book.Title,
			Author: analysis.Book.Authors[0].Name,
		}
		books = append(books, book)
	}
	return SuggestionResult{Books: books, BookAnalysis: bookAnalysis}, err
}
