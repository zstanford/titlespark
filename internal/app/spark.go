package app

import "fmt"

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

func SuggestBook(pref Preferences) (BookSpark, error) {
	openlibClient := NewOpenLibClient()
	qResult, err := openlibClient.OpenLib.Query(pref.Language, pref.Genre, pref.TargetAudience, pref.Subject)
	if err != nil {
		return BookSpark{}, err
	}
	bResult, err := openlibClient.OpenLib.GetBooks(qResult.Docs)
	if err != nil {
		fmt.Printf("26:%s", err)
	}
	book := BookSpark{
		Title:  bResult[0].Title,
		Author: bResult[0].Authors[0].Name,
	}
	fmt.Printf("Title: %s Author: %s", book.Title, book.Author)
	return book, err
}
