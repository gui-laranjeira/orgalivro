package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type ISBNResult struct {
	Title       string   `json:"title"`
	ISBN13      string   `json:"isbn13"`
	CoverURL    string   `json:"cover_url"`
	Description string   `json:"description"`
	Year        int      `json:"year"`
	Language    string   `json:"language"`
	Authors     []string `json:"authors"`
	Genres      []string `json:"genres"`
}

type ISBNService struct{ client *http.Client }

func NewISBNService() *ISBNService {
	return &ISBNService{client: &http.Client{}}
}

func (s *ISBNService) Lookup(isbn string) (*ISBNResult, error) {
	if result, err := s.lookupBrasilAPI(isbn); err == nil {
		return result, nil
	}
	return s.lookupOpenLibrary(isbn)
}

// lookupBrasilAPI queries brasilapi.com.br — excellent coverage for Brazilian ISBNs.
func (s *ISBNService) lookupBrasilAPI(isbn string) (*ISBNResult, error) {
	url := fmt.Sprintf("https://brasilapi.com.br/api/isbn/v1/%s", isbn)
	resp, err := s.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("brasilapi: status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data struct {
		Title    string   `json:"title"`
		Synopsis string   `json:"synopsis"`
		Authors  []string `json:"authors"`
		Subjects []string `json:"subjects"`
		Year     int      `json:"year"`
		CoverURL *string  `json:"cover_url"`
	}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}
	if data.Title == "" {
		return nil, fmt.Errorf("brasilapi: empty title")
	}

	result := &ISBNResult{
		ISBN13:      isbn,
		Title:       data.Title,
		Description: data.Synopsis,
		Year:        data.Year,
		Language:    "pt-BR",
		Authors:     data.Authors,
		Genres:      data.Subjects,
	}
	if data.CoverURL != nil {
		result.CoverURL = *data.CoverURL
	}
	// Fill cover from Open Library covers CDN if missing.
	if result.CoverURL == "" {
		result.CoverURL = fmt.Sprintf("https://covers.openlibrary.org/b/isbn/%s-L.jpg", isbn)
	}
	return result, nil
}

// lookupOpenLibrary queries openlibrary.org — fallback for non-Brazilian ISBNs.
func (s *ISBNService) lookupOpenLibrary(isbn string) (*ISBNResult, error) {
	url := fmt.Sprintf("https://openlibrary.org/api/books?bibkeys=ISBN:%s&format=json&jscmd=data", isbn)
	resp, err := s.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var raw map[string]any
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, err
	}

	key := "ISBN:" + isbn
	entryRaw, ok := raw[key]
	if !ok {
		return nil, fmt.Errorf("ISBN not found")
	}
	data, ok := entryRaw.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("unexpected response format")
	}

	result := &ISBNResult{ISBN13: isbn, Authors: []string{}, Genres: []string{}}

	if v, ok := data["title"].(string); ok {
		result.Title = v
	}
	if covers, ok := data["cover"].(map[string]any); ok {
		if large, ok := covers["large"].(string); ok {
			result.CoverURL = large
		} else if medium, ok := covers["medium"].(string); ok {
			result.CoverURL = medium
		}
	}
	if py, ok := data["publish_date"].(string); ok {
		for i := 0; i+4 <= len(py); i++ {
			var year int
			if n, _ := fmt.Sscanf(py[i:i+4], "%d", &year); n == 1 && year > 1000 && year < 2100 {
				result.Year = year
				break
			}
		}
	}
	if authorsRaw, ok := data["authors"].([]any); ok {
		for _, a := range authorsRaw {
			if am, ok := a.(map[string]any); ok {
				if name, ok := am["name"].(string); ok {
					result.Authors = append(result.Authors, name)
				}
			}
		}
	}
	if notes, ok := data["notes"].(map[string]any); ok {
		if v, ok := notes["value"].(string); ok {
			result.Description = v
		}
	} else if notes, ok := data["notes"].(string); ok {
		result.Description = notes
	}
	if subjects, ok := data["subjects"].([]any); ok {
		for _, sub := range subjects {
			if sm, ok := sub.(map[string]any); ok {
				if name, ok := sm["name"].(string); ok {
					result.Genres = append(result.Genres, strings.Title(strings.ToLower(name))) //nolint:staticcheck
				}
			}
		}
	}

	return result, nil
}
