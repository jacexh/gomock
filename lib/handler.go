package gomock

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/uber-go/zap"
)

type (
	mockResponse struct {
		Code    int    `json:"code"`
		Message string `json:"msg"`
	}
)

var (
	parseRuleFail = mockResponse{Code: 400, Message: "parse mock rule fail"}
	success       = mockResponse{Code: 200, Message: "success"}
)

// HandleMock for mocked api
func HandleMock(w http.ResponseWriter, r *http.Request) {
	Logger.Debug("got request", zap.String("url", r.URL.RequestURI()), zap.String("method", r.Method))

	id := genID(r.URL.Path, r.Method)
	rule := defaultPool.get(id)
	if rule == nil {
		rule = defaultRule
	}

	switch rule.Mode {
	case ModeNormal:
		for _, t := range rule.Templates {
			if t.IsMatched(r) {
				t.ToResponse(w)
				return
			}
		}
	case ModeKeyword:
		for _, t := range rule.Templates {
			if t.IsMatchedByKeyword(r) {
				t.ToResponse(w)
				return
			}
		}
	case ModeRegular:
		for _, t := range rule.Templates {
			if t.IsMatchedByRegular(r) {
				t.ToResponse(w)
				return
			}
		}
	}
}

// HandleCreate for create mock rule
func HandleCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		io.WriteString(w, "create a mock rule with POST method")
		return
	}
	defer r.Body.Close()
	rule := &MockRule{}
	err := json.NewDecoder(r.Body).Decode(rule)
	if err != nil {
		logError(err)
		resp, _ := json.Marshal(parseRuleFail)
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, string(resp))
		return
	}

	err = rule.filte()
	if err != nil {
		logError(err)
		ret := mockResponse{Code: 400, Message: err.Error()}
		resp, _ := json.Marshal(ret)
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, string(resp))
		return
	}

	defaultPool.receive(rule)
	w.Header().Set("Content-Type", "application/json")
	resp, _ := json.Marshal(success)
	io.WriteString(w, string(resp))
}

// HandleExport to export current mock rules
func HandleExport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		return
	}

	export := make([]*MockRule, len(defaultPool.pool))
	index := 0
	for _, v := range defaultPool.pool {
		export[index] = v
		index++
	}
	resp, err := json.MarshalIndent(export, "", "  ")
	if err != nil {
		logError(err)
		ret := mockResponse{Code: 500, Message: err.Error()}
		errResp, _ := json.Marshal(ret)

		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, string(errResp))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, string(resp))
}

// HandleImport to import local settings
func HandleImport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		return
	}

	rules := &[]*MockRule{}
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(rules)
	if err != nil {
		logError(err)
		ret := mockResponse{Code: 400, Message: err.Error()}
		errResp, _ := json.Marshal(ret)

		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, string(errResp))
		return
	}

	defaultPool.reset()
	defaultPool.batchReceive(*rules...)

	w.Header().Set("Content-Type", "application/json")
	resp, _ := json.Marshal(success)
	io.WriteString(w, string(resp))
}
