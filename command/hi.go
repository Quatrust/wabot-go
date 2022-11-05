package command

import (
	"fmt"
	"github.com/fdvky1/wabot-go/internal/handler"
	"github.com/fdvky1/wabot-go/util"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"time"
)

func HiCommand() {
	AddCommand(&handler.Command{
		Name:         "tes",
		Aliases:      []string{"hi"},
		Category:     handler.MiscCategory,
		RunFunc:      HiRunFunc,
		HideFromHelp: true,
	})
}

func HiRunFunc(c *whatsmeow.Client, args handler.RunFuncArgs) *waProto.Message {
	t := time.Now()
	util.SendReplyMessage(c, args.Evm, "testing a...")
	return util.SendReplyText(args.Evm, fmt.Sprintf("Duration: %f seconds", time.Now().Sub(t).Seconds()))
}
