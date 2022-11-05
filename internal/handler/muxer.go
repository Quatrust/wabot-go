package handler

import (
	"fmt"
	"github.com/fdvky1/wabot-go/helper"
	"github.com/fdvky1/wabot-go/util"
	"github.com/zhangyunhao116/skipmap"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types"
	"golang.org/x/exp/slices"
	"strings"
)

var GlobalMiddleware, MsgMiddleware *skipmap.StringMap[MiddlewareFunc]

type Muxer struct {
	CmdMap *skipmap.StringMap[*Command]
	Locals *skipmap.StringMap[string]
}

func (m *Muxer) AddCommand(cmd *Command) {
	cmd.Validate()
	_, ok := m.CmdMap.Load(cmd.Name)
	if ok {
		panic("Duplicate command: " + cmd.Name)
	}

	for _, alias := range cmd.Aliases {
		_, ok := m.CmdMap.Load(alias)
		if ok {
			panic("Duplicate alias in command " + cmd.Name)
		}
		m.CmdMap.Store(alias, cmd)
	}
	m.CmdMap.Store(cmd.Name, cmd)
}

func (m *Muxer) CheckGlobalState(number string) (bool, string) {
	globalState, ok := m.Locals.Load(number)
	if !ok {
		return false, ""
	}
	m.Locals.Delete(number)
	return true, globalState
}

func (m *Muxer) RunCommand(c *whatsmeow.Client, args RunFuncArgs) {
	ok, stateCmd := m.CheckGlobalState(args.Number)
	if ok {
		if strings.Contains(args.Msg, "!cancel") {
			m.Locals.Delete(args.Number)
			util.SendReplyMessage(c, args.Evm, "Canceled!")
			return
		}
		args.Msg = stateCmd + " " + args.Msg
	}
	
	if args.Cmd != nil {
		GlobalMiddleware.Range(func(key string, value MiddlewareFunc) bool {
			if !value(c, args) {
				return false
			}
			return true
		})

		if args.Cmd.Middleware != nil {
			if !args.Cmd.Middleware(c, args) {
				return
			}
		}
		if args.Cmd.GroupOnly {
			if !args.Evm.Info.IsGroup {
			  util.SendReplyMessage(c, args.Evm, "Sorry, commands can only be used within a group!")
				return
			}
		}
		if args.Cmd.AdminOnly && args.Evm.Info.IsGroup {
		  GroupInfo, err := c.GetGroupInfo(args.Evm.Info.Chat)
		  if err != nil {
		    return
		  }
		  user := GroupInfo.Participants[slices.IndexFunc(GroupInfo.Participants, func(c types.GroupParticipant) bool { return c.JID == args.Evm.Info.Sender})]
		  if !user.IsAdmin {
		    util.SendReplyMessage(c, args.Evm, "Sorry, Only admin can use this command!")
				return
		  }
		}
		if args.Cmd.PrivateOnly {
			if args.Evm.Info.IsGroup {
			  util.SendReplyMessage(c, args.Evm, "Sorry, commands can only be used within a private message!")
				return
			}
		}
		msg := args.Cmd.RunFunc(c, args)

		if msg != nil {
			_, err := c.SendMessage(args.Evm.Info.Chat, "", msg)
			if err != nil {
				fmt.Println(err)
			}
		}
		
		/*
		a := []types.MessageID{}
		a = append(a, args.Evm.Info.ID)
		c.MarkRead(a, args.Evm.Info.Timestamp, args.Evm.Info.Chat, args.Evm.Info.Sender)
		*/
	}
}

func (m *Muxer) GenerateRequiredLocals() {
	uid := helper.CreateUid()
	m.Locals.Store("uid", uid)
}

func NewMuxer() *Muxer {
	muxer := &Muxer{
		Locals: skipmap.NewString[string](),
		CmdMap: skipmap.NewString[*Command](),
	}
	muxer.GenerateRequiredLocals()
	return muxer
}
