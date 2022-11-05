package MessageMiddleware

import (
  "fmt"
	"github.com/fdvky1/wabot-go/internal/handler"
	//"github.com/fdvky1/wabot-go/util"
	"go.mau.fi/whatsmeow"
)

func LogMiddleware(c *whatsmeow.Client, args handler.RunFuncArgs) bool {
	if args.Cmd == nil && args.Msg != "" {
	  fmt.Println(fmt.Sprintf("[MSG] [%s] message : %s", args.Number, args.Msg))
	}
	return true
}
