package MessageMiddleware

import (
  "os"
	"github.com/fdvky1/wabot-go/internal/handler"
	"go.mau.fi/whatsmeow"
)

func SelfModeMiddleware(c *whatsmeow.Client, args handler.RunFuncArgs) bool {
	if os.Getenv("MODE") == "self" && !args.Evm.Info.MessageSource.IsFromMe {
	  return false
	}
	return true
}

