package MessageMiddleware

import (
  "fmt"
  "strings"
  "net/http"
	"github.com/fdvky1/wabot-go/internal/handler"
	"github.com/fdvky1/wabot-go/util"
	"go.mau.fi/whatsmeow"
	"github.com/gocolly/colly/v2"
)

func AxisBuyMiddleware(c *whatsmeow.Client, args handler.RunFuncArgs) bool {
  QuotedMsg := util.ParseQuotedMessage(args.Evm.Message)
  if QuotedMsg != nil && QuotedMsg.ListMessage != nil {
    cookie := fmt.Sprintf("ci_session=%s", QuotedMsg.ListMessage.GetFooterText())
    co := colly.NewCollector()
    co_ := co.Clone()
    co__ := co.Clone()
    co.OnRequest(func(r *colly.Request) {
      r.Headers.Set("cookie", cookie)
    })
    co_.OnRequest(func(r *colly.Request) {
      r.Headers.Set("cookie", cookie)
    })
    co__.OnRequest(func(r *colly.Request) {
      r.Headers.Set("cookie", cookie)
    })
    co.OnHTML("table#table-datatable > tbody", func(h *colly.HTMLElement) {
      h.ForEach("tr", func(_ int, el *colly.HTMLElement) {
        if el.ChildText("td:nth-child(2)") == args.Evm.Message.ListResponseMessage.GetTitle() {
          co_.Visit(el.ChildAttr("td:nth-child(4) a", "href"))
        }
      })
    })
    co_.OnHTML(`form[action="https://indo5.com/Axis/buy_paket"]`, func(h *colly.HTMLElement) {
      co__.Request("POST", "https://indo5.com/Axis/buy_paket", strings.NewReader(fmt.Sprintf("type=%s&paket=%s", h.ChildAttr("#type", "value"), h.ChildAttr("#paket", "value"))), nil, http.Header{"Content-Type": []string{"application/x-www-form-urlencoded"}})
    })
    co__.OnHTML(`form[action="https://indo5.com/Axis/dashboard"]`, func(h *colly.HTMLElement) {
      text := fmt.Sprintf("Trying to buy an %s package, for more info please check the message from the operator", args.Evm.Message.ListResponseMessage.GetTitle())
      if h.ChildText(".alert-warning") == "Maaf, paket tidak tersedia" {
        text = "Sorry, package is unavailable"
      }
      util.SendText(c, args.Evm, text)
    })
    co.Visit("https://indo5.com/Axis/dashboard")
  }
	return true
}
