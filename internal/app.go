package internal

import (
  "strings"
	"context"
	"fmt"
	"github.com/fdvky1/wabot-go/command"
	"github.com/fdvky1/wabot-go/config"
	"os"
	"os/signal"
	"syscall"
	"time"
	"strconv"

	"github.com/fdvky1/wabot-go/internal/handler"
	"github.com/fdvky1/wabot-go/util"
	"github.com/mdp/qrterminal/v3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types/events"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

type App struct {
	Log      *zap.Logger
	SqlStore *sqlstore.Container
}

type Base struct {
	startTime time.Time
	Log       *zap.Logger
	Muxer     *handler.Muxer
	client    *whatsmeow.Client
	Device    *store.Device
}

func (b *Base) QRChanFunc(ch <-chan whatsmeow.QRChannelItem) {
	for evt := range ch {
		if evt.Event == "code" {
			qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
		} else {
			b.Log.Info("QRChan", zap.Any("QR channel result", evt.Event))
		}
	}
}

func (b *Base) ConnectedEvents(evt interface{}) {
	_, ok := evt.(*events.Connected)
	if ok {
		b.startTime = time.Now()
		fmt.Println("[INFO] Connected!")
	}
}

func (b *Base) MessageEvents(evt interface{}) {
	event, ok := evt.(*events.Message)
	if ok {
		if !b.startTime.IsZero() && event.Info.Timestamp.After(b.startTime) {
		  id, _ := b.Muxer.Locals.Load("uid")
		  parsed := util.ParseMessageText(id, event)
		  cmd, _ := util.ParseCmd(parsed)
		  command, _ := b.Muxer.CmdMap.Load(cmd)
	    args := handler.RunFuncArgs{
  			Evm:    event,
  			Cmd:    command,
  			Msg:    parsed,
  			Number: strconv.FormatUint(event.Info.Sender.UserInt(), 10),
  			Locals: b.Muxer.Locals,
  			Args:   strings.Split(parsed, " "),
  		}
  		
  		ContinueByMiddleware := false
		  handler.MsgMiddleware.Range(func(key string, value handler.MiddlewareFunc) bool {
  			if !value(b.client, args) {
  			  ContinueByMiddleware = true
  			}
  			return true
		  })
		  if ContinueByMiddleware {
		    return
		  }
		 	go b.Muxer.RunCommand(b.client, args)
		}
	}
}

func (b *Base) Init() {
	store.DeviceProps.RequireFullSync = proto.Bool(true)
	b.Muxer = handler.NewMuxer()

	for _, cmd := range command.Commands {
		b.Muxer.AddCommand(cmd)
	}

	b.client = whatsmeow.NewClient(b.Device, config.NewWaLog())
	if b.client.Store.ID == nil {
		qrChan, _ := b.client.GetQRChannel(context.Background())
		err := b.client.Connect()
		if err != nil {
			b.Log.With(zap.Error(err)).Error(err.Error())
		}
		go b.QRChanFunc(qrChan)
	} else {
		err := b.client.Connect()
		if err != nil {
			b.Log.With(zap.Error(err)).Error(err.Error())
		}
	}
	b.client.AddEventHandler(b.ConnectedEvents)
	b.client.AddEventHandler(b.MessageEvents)

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	b.client.Disconnect()
}
