package gomock

import (
	"runtime"

	"strings"

	"github.com/uber-go/zap"
)

// Logger global logger
var Logger zap.Logger

func init() {
	Logger = zap.New(zap.NewJSONEncoder(zap.RFC3339Formatter("time"))).With(zap.String("package", "gomock"))

	defaultPool = MockRulePool{pool: map[string]*MockRule{}}
	defaultRule = &MockRule{
		Path:   "/",
		Method: "GET",
		Mode:   ModeNormal,
		Templates: []*Template{
			&Template{
				Content:     `{"server": "gomock"}`,
				StatusCode:  200,
				ContentType: "application/json",
			},
		},
	}
}

func logError(e error) {
	if e != nil {
		pc, _, line, _ := runtime.Caller(1)
		Logger.Error(e.Error(),
			zap.String("func", runtime.FuncForPC(pc).Name()),
			zap.Int("line", line))
	}
}

func genID(path, method string) string {
	return path + "||" + strings.ToUpper(method)
}
