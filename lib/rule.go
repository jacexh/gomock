package gomock

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
)

type (
	// MockRule the rule of mocked api
	MockRule struct {
		Path      string      `json:"path"`
		Method    string      `json:"method"`
		Mode      int         `json:"mode"`
		Templates []*Template `json:"responses"`
	}
	// MockRulePool collection of MockRule
	MockRulePool struct {
		pool map[string]*MockRule
	}
)

var (
	defaultPool MockRulePool
	defaultRule *MockRule
)

func (r *MockRule) id() string {
	return genID(r.Path, r.Method)
}

func (r *MockRule) filte() error {
	r.Method = strings.ToUpper(r.Method)

	if len(r.Templates) == 0 {
		return errors.New("lack of requeired field: responses")
	}

	switch r.Mode {
	case ModeKeyword:
		for _, t := range r.Templates {
			if len(t.Keyword) == 0 {
				return errors.New("lack of required field: keyword")
			}
		}
	case ModeRegular:
		for _, t := range r.Templates {
			if len(t.Regular) == 0 {
				return errors.New("lack of requried field: keyword")
			}
		}
	}
	return nil
}

// Get get mock rule from rule pool by id
func (p *MockRulePool) Get(id string) *MockRule {
	rule, ok := p.pool[id]
	if ok {
		return rule
	}
	return nil
}

// Receive insert MockRule to MockRulePool
func (p *MockRulePool) Receive(r *MockRule) {
	p.pool[r.id()] = r
}

func (p *MockRulePool) handleExport(w http.ResponseWriter, r *http.Request) {
	export := make([]*MockRule, len(p.pool))
	index := 0
	for _, v := range p.pool {
		export[index] = v
		index++
	}
	resp, err := json.Marshal(export)
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
