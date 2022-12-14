package command

import (
  "fmt"
  "strings"
	"github.com/fdvky1/wabot-go/internal/handler"
	"google.golang.org/protobuf/proto"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	
)

func EveryoneCommand() {
	AddCommand(&handler.Command{
		Name:        "everyone",
		Aliases:     []string{"hidetag", "ht", "tagall", "mentionall"},
		Category:    handler.GroupCategory,
		RunFunc:     EveryoneRunFunc,
		GroupOnly:   true,
		AdminOnly:   true,
	})
}

func EveryoneRunFunc(c *whatsmeow.Client, args handler.RunFuncArgs) *waProto.Message {
  info, err := c.GetGroupInfo(args.Evm.Info.Chat)
  if err != nil {
    return nil
  }
  member := []string{}
  for _, m := range info.Participants {
    member = append(member, m.JID.String())
  }
  text := "@everyone"
  if len(args.Args) > 1 {
    text = fmt.Sprintf("%s\n\n%s", strings.Replace(args.Msg, args.Args[0], "", 1), text)
  }
  msg :=  &waProto.Message{
		ExtendedTextMessage: &waProto.ExtendedTextMessage{
			Text:        proto.String(text),
			ContextInfo: &waProto.ContextInfo {
			  MentionedJid: member,
			},
		},
	}
	return msg
}
