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

func SuggestBook(pref Preferences) ([]BookSpark, error) {
	var books []BookSpark
	openlibClient := NewOpenLibClient()
	qResult, err := openlibClient.OpenLib.Query(pref.Language, pref.Genre, pref.TargetAudience, pref.Subject)
	if err != nil {
		fmt.Printf("Error during query: ", err)
	}
	bResult, err := openlibClient.OpenLib.GetBooks(qResult.Docs)
	if err != nil {
		fmt.Printf("Error during get books: ", err)
	}

	if len(bResult) == 0 {
		return []BookSpark{}, errors.New("no books found")
	}
	fmt.Printf("Book Result length: %v\n", len(bResult))
	for _, result := range bResult {
		book := BookSpark{
			Title:  result.Title,
			Author: result.Authors[0].Name,
		}
		books = append(books, book)
	}
	return books, err
}
