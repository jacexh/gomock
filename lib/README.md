# gomock
--
    import "github.com/jacexh/gomock/lib"


## Usage

```go
const (
	// ModeNormal 0
	ModeNormal = iota
	// ModeKeyword 1
	ModeKeyword
	// ModeRegular 2
	ModeRegular
)
```

```go
var Logger zap.Logger
```
Logger global logger

#### func  HandleCreate

```go
func HandleCreate(w http.ResponseWriter, r *http.Request)
```
HandleCreate for create mock rule

#### func  HandleExport

```go
func HandleExport(w http.ResponseWriter, r *http.Request)
```
HandleExport to export current mock rules

#### func  HandleImport

```go
func HandleImport(w http.ResponseWriter, r *http.Request)
```
HandleImport to import local settings

#### func  HandleMock

```go
func HandleMock(w http.ResponseWriter, r *http.Request)
```
HandleMock for mocked api

#### type MockRule

```go
type MockRule struct {
	Path      string      `json:"path"`
	Method    string      `json:"method"`
	Mode      int         `json:"mode"`
	Templates []*Template `json:"responses"`
}
```

MockRule the rule of mocked api

#### type MockRulePool

```go
type MockRulePool struct {
}
```

MockRulePool collection of MockRule

#### type Template

```go
type Template struct {
	Content     string            `json:"content"`
	ContentType string            `json:"content_type,omitempty"`
	StatusCode  int               `json:"status_code,omitempty"`
	Headers     map[string]string `json:"headers,omitempty"`
	Keyword     string            `json:"keyword"`
	Regular     string            `json:"regular"`
}
```

Template return fixed response

#### func (*Template) IsMatched

```go
func (t *Template) IsMatched(r *http.Request) bool
```
IsMatched always return true

#### func (*Template) IsMatchedByKeyword

```go
func (t *Template) IsMatchedByKeyword(r *http.Request) bool
```
IsMatchedByKeyword defect the http.request if contains specific keyword

#### func (*Template) IsMatchedByRegular

```go
func (t *Template) IsMatchedByRegular(r *http.Request) bool
```
IsMatchedByRegular defect the http request if match the regular option

#### func (*Template) ToResponse

```go
func (t *Template) ToResponse(w http.ResponseWriter)
```
ToResponse render the template to http response
