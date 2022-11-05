package MessageMiddleware

import (
  "fmt"
  "net/http"
  "strings"
	"github.com/fdvky1/wabot-go/internal/handler"
	"github.com/fdvky1/wabot-go/util"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"github.com/gocolly/colly/v2"
)

func AxisOtpMiddleware(c *whatsmeow.Client, args handler.RunFuncArgs) bool {
	quotedMsg := util.ParseQuotedMessage(args.Evm.Message).GetConversation()
	if strings.Contains(quotedMsg, "AXIS#") && args.Msg != "" {
	  cookie := fmt.Sprintf("ci_session=%s", strings.Split(quotedMsg, "AXIS#")[1])
	  client := &http.Client{}
	  req, _ := http.NewRequest("POST", "https://indo5.com/Axis/validate_otp", strings.NewReader(fmt.Sprintf("otp=%s", strings.ToUpper(args.Msg))))
    req.Header.Set("Cookie", cookie)
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
    _, err := client.Do(req)
    if err != nil {
      util.SendReplyMessage(c, args.Evm, "Failed to validate otp, try request otp again")
    }
    co := colly.NewCollector()
    SectionRow := []*waProto.Row{}
    co.OnRequest(func(r *colly.Request) {
      r.Headers.Set("cookie", cookie)
    })
    co.OnHTML("table#table-datatable > tbody", func(h *colly.HTMLElement) {
      h.ForEach("tr", func(_ int, el *colly.HTMLElement) {
        SectionRow = append(SectionRow, util.CreateSectionRow(el.ChildText("td:nth-child(2)"), fmt.Sprintf("Rp. %s", el.ChildText("td:nth-child(3)")), el.ChildText("td:nth-child(1)")))
      })
      msg := util.GenerateListMessage("Axis Package Price List", "You can buy a data package here\n\nNote:\n*- Make sure you have enough credit*\n*- Use at your own risk!*", "select", strings.Split(cookie, "=")[1], util.CreateSectionList("Choose a package to buy", SectionRow))
      _, err := c.SendMessage(args.Evm.Info.Chat, "", msg)
			if err != nil {
				fmt.Println(err)
			}
    })
    co.Visit("https://indo5.com/Axis/dashboard")
	}
	return true
}