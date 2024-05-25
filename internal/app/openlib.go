package app

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
	"unicode"
)

const openlibUrl = "https://openlibrary.org"

type OpenLibClient struct {
	Client  *http.Client
	baseUrl *url.URL

	OpenLib OpenLibService
}

func NewOpenLibClient() *OpenLibClient {
	u, _ := url.Parse(openlibUrl)
	c := &OpenLibClient{
		Client: &http.Client{
			Timeout: time.Second * 100,
		},
		baseUrl: u,
	}
	c.OpenLib = &OpenLibServiceOp{client: c}
	return c
}

func (c *OpenLibClient) NewRequest(path string) (*http.Response, error) {
	u, _ := url.Parse(c.baseUrl.String() + path)
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	fmt.Printf("Req: %v", req)
	res, err := c.Client.Do(req)
	if err != nil {
		fmt.Println(err)
		return res, err
	}
	return res, err
}

type QueryResponse struct {
	NumFound      int        `json:"numFound"`
	Start         int        `json:"start"`
	NumFoundExact bool       `json:"numFoundExact"`
	Docs          []BookISBN `json:"docs"`
}

type BookISBN struct {
	ISBN []string `json:"isbn"`
}

type OpenLibService interface {
	Query(string, string, string, string) (*QueryResponse, error)
	GetBooks([]BookISBN) ([]Book, error)
}

type OpenLibServiceOp struct {
	client *OpenLibClient
}

type Book struct {
	Title       string    `json:"title"`
	Authors     []Author  `json:"authors"`
	Subjects    []Subject `json:"subjects"`
	PublishDate string    `json:"publish_date"`
	Pages       int       `json:"number_of_pages"`
}

type Author struct {
	Name string `json:"name"`
	Url  string `json:url`
}

type Subject struct {
	Name string `json:"name"`
	Url  string `json:url`
}

type SearchResults struct {
	ISBN []string `json:"isbn"`
}

func (s *OpenLibServiceOp) Query(language string, genre string, age string, subject string) (*QueryResponse, error) {
	q := fmt.Sprintf("subject:%s subject:%s subject:%s language:%s", genre, age, subject, language)
	encodedQ := url.QueryEscape(q)
	fields := "isbn"
	limit := "100"
	queryUrl := fmt.Sprintf("/search.json?q=%s&fields=%s&limit=%s", encodedQ, fields, limit)
	res, err := s.client.NewRequest(queryUrl)
	fmt.Printf("res:%v", res)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	data, _ := io.ReadAll(res.Body)
	var result QueryResponse
	err = json.Unmarshal(data, &result)
	return &result, err
}

// TODO: still getting error
func (s *OpenLibServiceOp) GetBooks(booksISBN []BookISBN) ([]Book, error) {
	var (
		books     []Book
		bibKeyStr string
	)

	if len(booksISBN) == 0 {
		return nil, nil
	}

	for _, book := range booksISBN {
		for _, isbn := range book.ISBN {
			firstISBN := extractFirstSetOfNumbers(isbn)
			if firstISBN == "" {
				continue
			}
			bibKeyStr += fmt.Sprintf("ISBN:%s,", firstISBN)
		}
	}
	queryUrl := fmt.Sprintf("/api/books?bibkey=%s&jscmd=data&format=json", bibKeyStr)
	res, err := s.client.NewRequest(queryUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	data, _ := io.ReadAll(res.Body)
	bookMap := make(map[string]Book)
	err = json.Unmarshal([]byte(data), &bookMap)
	if err != nil {
		return nil, err
	}
	for _, book := range bookMap {
		books = append(books, book)
	}
	return books, nil
}

func extractFirstSetOfNumbers(s string) string {
	var numbers []rune
	for _, r := range s {
		if unicode.IsDigit(r) {
			numbers = append(numbers, r)
		} else if len(numbers) > 0 {
			break
		}
	}
	return string(numbers)
}
