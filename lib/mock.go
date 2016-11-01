package gomock

import (
	"runtime"

	"github.com/uber-go/zap"
)

// Logger global logger
var Logger zap.Logger

func init() {
	Logger = zap.New(zap.NewJSONEncoder(zap.RFC3339Formatter("time"))).With(zap.String("package", "gomock"))
}

func logError(e error) {
	if e != nil {
		pc, _, line, _ := runtime.Caller(1)
		Logger.Error(e.Error(),
			zap.String("func", runtime.FuncForPC(pc).Name()),
			zap.Int("line", line))
	}
}
