package app

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
	"titlespark-web/internal/util"
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
	res, err := c.Client.Do(req)
	if err != nil {
		return res, err
	}
	return res, err
}

type QueryResponse struct {
	NumFound      int        `json:"numFound"`
	Start         int        `json:"start"`
	NumFoundExact bool       `json:"numFoundExact"`
	Docs          []BookDocs `json:"docs"`
}

type GetBookResponse struct {
	Title   string   `json:"title"`
	Authors []Author `json:"authors"`
}

type BookDocs struct {
	ISBN     []string `json:"isbn"`
	Subjects []string `json:"subject"`
}

type OpenLibService interface {
	Query(string, string, string, string, int) (*QueryResponse, error)
	GetBooks([]BookDocs) ([]Book, error)
}

type OpenLibServiceOp struct {
	client *OpenLibClient
}

type Book struct {
	Title    string
	Author   string
	Subjects []string
}

type Author struct {
	Name string `json:"name"`
	Url  string `json:url`
}

func (s *OpenLibServiceOp) Query(language string, genre string, age string, subject string, limit int) (*QueryResponse, error) {
	sanitizedSubject := util.RemoveSpecialCharacters(subject)
	if util.IsOnlySpaces(sanitizedSubject) || sanitizedSubject == "" {
		sanitizedSubject = "general" // general subject will return more results instead of empty string
	}
	q := fmt.Sprintf("subject:%s subject:%s language:%s", genre, sanitizedSubject, language)
	encodedQ := url.QueryEscape(q)
	fields := "isbn,subject"
	queryUrl := fmt.Sprintf("/search.json?q=%s&fields=%s&limit=%v", encodedQ, fields, limit)
	res, err := s.client.NewRequest(queryUrl)
	if err != nil {
		return nil, err
	}
	resContentType := res.Header.Get("Content-Type")
	if strings.Contains(resContentType, "text/html") {
		return &QueryResponse{}, err
	}
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("error reading body: ", err)
	}
	var result QueryResponse
	err = json.Unmarshal(data, &result)
	if err != nil {
		fmt.Printf("error marshalling: ", err)
	}
	return &result, err
}

func (s *OpenLibServiceOp) GetBooks(bookDocs []BookDocs) ([]Book, error) {
	var (
		books     []Book
		bibKeyStr string
	)

	if len(bookDocs) == 0 {
		return nil, fmt.Errorf("no ISBNs")
	}

	isbnToSubjects := make(map[string][]string)
	for _, book := range bookDocs {
		if len(book.ISBN) > 0 {
			firstISBN := book.ISBN[0]
			isbnToSubjects[firstISBN] = book.Subjects
			bibKeyStr += fmt.Sprintf("ISBN:%s,", firstISBN)
		}
	}
	queryUrl := fmt.Sprintf("/api/books?bibkeys=%s&jscmd=data&format=json", bibKeyStr)
	res, err := s.client.NewRequest(queryUrl)
	if err != nil {
		return []Book{}, err
	}
	resContentType := res.Header.Get("Content-Type")
	if strings.Contains(resContentType, "text/html") {
		return []Book{}, err
	}
	defer res.Body.Close()
	data, _ := io.ReadAll(res.Body)
	bookMap := make(map[string]GetBookResponse)
	err = json.Unmarshal([]byte(data), &bookMap)
	if err != nil {
		return []Book{}, err
	}
	for bibKey, book := range bookMap {
		isbn := strings.TrimPrefix(bibKey, "ISBN:")
		subjects, exists := isbnToSubjects[isbn]
		if !exists {
			return []Book{}, fmt.Errorf("subjects not found for ISBN: %s", isbn)
		}
		b := Book{
			Title:    book.Title,
			Author:   book.Authors[0].Name,
			Subjects: subjects,
		}
		books = append(books, b)
	}
	return books, nil
}
