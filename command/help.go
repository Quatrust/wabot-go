package command

import (
	"fmt"
	"time"
	"strings"
	"github.com/fdvky1/wabot-go/internal/handler"
	"github.com/fdvky1/wabot-go/util"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
)

func HelpCommand() {
	AddCommand(&handler.Command{
		Name:        "help",
		Aliases:     []string{"menu","h"},
		Category:    handler.UtilitiesCategory,
		RunFunc:     HelpRunFunc,
	})
}

func generateCmdText(cmdList []*handler.Command, prefix string) string {
  text := fmt.Sprintf("*📃%s*\n", cmdList[0].Category.Name)
  for i, v := range cmdList {
    if i == (len(cmdList) - 1) {
      text += fmt.Sprintf("*└ %s%s*\n\n", prefix, v.Name)
    } else {
      text += fmt.Sprintf("*├ %s%s*\n", prefix, v.Name)
    }
  }
  return text
}

func HelpRunFunc(c *whatsmeow.Client, args handler.RunFuncArgs) *waProto.Message {
  prefix := strings.Split(args.Args[0], "")[0]
  t := time.Now()
  hour := t.Hour() 
  greetings := "Konbanwa"
  if hour > 01 && hour < 12 {
    greetings = "Ohayō"
  } else if hour > 12 && hour < 18 {
    greetings = "Konnichiwa"
  }
  
  var Utilities, Misc []*handler.Command
  for _, data := range Commands {
    if data.HideFromHelp {
      continue
    }
    if data.Category.Name == "Utilities"{
      Utilities = append(Utilities, data)
    } else if data.Category.Name == "Misc"{
      Misc = append(Misc, data)
    }
  }
  text := fmt.Sprintf("%s *%s*👋\n📬 Need help? Here are all of my commands:\n\n", greetings, args.Evm.Info.PushName)
  text += generateCmdText(Utilities, prefix)
  text += generateCmdText(Misc, prefix)
	return util.SendReplyText(args.Evm, text)
}
