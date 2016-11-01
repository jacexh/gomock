package gomock

type (
	// MockRule the rule of mocked api
	MockRule struct {
		Path      string `json:"path"`
		Method    string `json:"method"`
		Mode      int    `json:"mode"`
		Templates []Render
		Responses []byte `json:"response"`
	}
	// MockRulePool collection of MockRule
	MockRulePool struct {
		pool map[string]*MockRule
	}
)

var (
	defaultPool MockRulePool
	defaultRule = MockRule{
		Path:   "/",
		Method: "GET",
		Mode:   ModeNormal,
		Templates: []Render{
			&TemplateBase{
				Content:     "Welcome to gomock",
				StatusCode:  200,
				ContentType: "text/html",
			},
		},
	}
)

func (r *MockRule) id() string {
	return r.Path + "||" + r.Method
}

// Get get mock rule from rule pool by id
func (p MockRulePool) Get(id string) *MockRule {
	rule, ok := p.pool[id]
	if ok {
		return rule
	}
	return nil
}

// Receive insert MockRule to MockRulePool
func (p MockRulePool) Receive(r *MockRule) {
	p.pool[r.id()] = r
}
