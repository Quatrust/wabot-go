package command

import (
	"fmt"
	"net/http"
	"net/url"
	"github.com/fdvky1/wabot-go/internal/handler"
	"github.com/fdvky1/wabot-go/util"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
)

func AxisCommand() {
	AddCommand(&handler.Command{
		Name:        "axis",
		Aliases:     []string{"tembak"},
		Category:    handler.MiscCategory,
		RunFunc:     AxisRunFunc,
		PrivateOnly: true,
	})
}

func AxisRunFunc(c *whatsmeow.Client, args handler.RunFuncArgs) *waProto.Message {
  if len(args.Args) < 2 {
    return util.SendReplyText(args.Evm, "Invalid Number")
  }
  form := url.Values{}
  form.Add("msisdn", args.Args[1])
  res, err := http.PostForm("https://indo5.com/Axis/validate", form)
  if err != nil{
    fmt.Printf("Failed to request otp: %v\n", err)
  }
	return util.SendReplyText(args.Evm, fmt.Sprintf("Otp sent successfully, reply to this message with code otp\n\nAXIS#%s",res.Cookies()[0].Value))
}
