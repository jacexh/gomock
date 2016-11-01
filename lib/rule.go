package gomock

import (
	"errors"
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

func (p *MockRulePool) get(id string) *MockRule {
	rule, ok := p.pool[id]
	if ok {
		return rule
	}
	return nil
}

func (p *MockRulePool) receive(r *MockRule) {
	p.pool[r.id()] = r
}

func (p *MockRulePool) reset() {
	p.pool = map[string]*MockRule{}
}

func (p *MockRulePool) batchReceive(rules ...*MockRule) {
	for _, r := range rules {
		p.receive(r)
	}
}
