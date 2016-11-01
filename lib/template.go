package gomock

import (
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

const (
	// ModeNormal 0
	ModeNormal = iota
	// ModeKeyword 1
	ModeKeyword
	// ModeRegular 2
	ModeRegular
)

type (
	// TemplateBase return fixed response
	TemplateBase struct {
		Content     string            `json:"content"`
		ContentType string            `json:"content_type,omitempty"`
		StatusCode  int               `json:"status_code,omitempty"`
		Headers     map[string]string `json:"headers,omitempty"`
		Keyword     string            `json:"keyword"`
		Regular     string            `json:"regular"`
	}

	//TemplateKeyword return by keyword
	TemplateKeyword struct {
		*TemplateBase
		Keyword string `json:"keyword"`
	}

	// TemplateRegular return by regular
	TemplateRegular struct {
		*TemplateBase
		Regular string `json:"regular"`
	}

	// Render .
	Render interface {
		IsMatched(*http.Request) bool
		ToResponse(http.ResponseWriter)
	}
)

// IsMatched always return true
func (t *TemplateBase) IsMatched(r *http.Request) bool {
	return true
}

// ToResponse render the template to http response
func (t *TemplateBase) ToResponse(w http.ResponseWriter) {
	if t.StatusCode != 0 {
		w.WriteHeader(t.StatusCode)
	}
	if t.ContentType != "" {
		w.Header().Set("Content-Type", t.ContentType)
	}
	for k, v := range t.Headers {
		w.Header().Set(k, v)
	}
	io.WriteString(w, t.Content)
}

// IsMatched defect the http.request if contains specific keyword
func (t *TemplateKeyword) IsMatched(r *http.Request) bool {
	if r.Method == http.MethodGet {
		return strings.Contains(r.URL.RawQuery, t.Keyword)
	}
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logError(err)
		return false
	}
	return strings.Contains(string(body), t.Keyword)
}

// IsMatched defect the http request if match the regular option
func (t *TemplateRegular) IsMatched(r *http.Request) bool {
	c := regexp.MustCompile(t.Regular)

	if r.Method == http.MethodGet {
		return c.MatchString(r.URL.RawQuery)
	}
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logError(err)
		return false
	}
	return c.MatchString(string(body))
}
