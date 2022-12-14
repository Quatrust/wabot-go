package command

import (
	"github.com/fdvky1/wabot-go/internal/handler"
	"github.com/fdvky1/wabot-go/util"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
)

func AddParticipantsCommand() {
	AddCommand(&handler.Command{
		Name:         "add",
		Aliases:      []string{"invite"},
		Category:     handler.GroupCategory,
		RunFunc:      AddParticipantsRunFunc,
		GroupOnly:    true,
		AdminOnly:    true,
	})
}

func AddParticipantsRunFunc(c *whatsmeow.Client, args handler.RunFuncArgs) *waProto.Message {
  if len(args.Args) < 2 {
    return util.SendReplyText(args.Evm, "Invalid number\nEx: /add 628...")
  }
  var Members map[types.JID]whatsmeow.ParticipantChange
  for _, v := range args.Args[1:len(args.Args)] {
    jid, ok := util.ParseJID(v)
    if ok {
      found, err := c.IsOnWhatsApp([]string{jid.User})
      if err == nil && len(found) != 0 {
        Members = make(map[types.JID]whatsmeow.ParticipantChange)
        Members[jid] = whatsmeow.ParticipantChangeAdd
        c.UpdateGroupParticipants(args.Evm.Info.Chat, Members)
      }
    }
  }
  return nil
}