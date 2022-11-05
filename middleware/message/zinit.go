package MessageMiddleware

import (
	"github.com/fdvky1/wabot-go/helper"
	"github.com/fdvky1/wabot-go/internal/handler"
)

func GenerateAllMiddlewares() {
	AddMiddleware(AxisOtpMiddleware)
	AddMiddleware(AxisBuyMiddleware)
	AddMiddleware(LogMiddleware)
	AddMiddleware(SelfModeMiddleware)
}

func AddMiddleware(mid handler.MiddlewareFunc) {
	handler.MsgMiddleware.Store(helper.CreateUid(), mid)
}
