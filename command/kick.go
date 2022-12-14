package command

import (
	"github.com/fdvky1/wabot-go/internal/handler"
	"github.com/fdvky1/wabot-go/util"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"strings"
)

func KickParticipantsCommand() {
	AddCommand(&handler.Command{
		Name:         "kick",
		Aliases:      []string{"remove"},
		Category:     handler.GroupCategory,
		RunFunc:      KickParticipantsRunFunc,
		GroupOnly:    true,
		AdminOnly:    true,
	})
}

func KickParticipantsRunFunc(c *whatsmeow.Client, args handler.RunFuncArgs) *waProto.Message {
  if len(args.Args) < 2 {
    return util.SendReplyText(args.Evm, "Mention member\nEx: /kick @member")
  }
  var Members map[types.JID]whatsmeow.ParticipantChange
  for _, v := range args.Args[1:len(args.Args)] {
    jid, ok := util.ParseJID(strings.Replace(v, "@", "", 1))
    if ok {
      found, err := c.IsOnWhatsApp([]string{jid.User})
      if err == nil && len(found) != 0 {
        Members = make(map[types.JID]whatsmeow.ParticipantChange)
        Members[jid] = whatsmeow.ParticipantChangeRemove
        c.UpdateGroupParticipants(args.Evm.Info.Chat, Members)
      }
    }
  }
  return nil
}