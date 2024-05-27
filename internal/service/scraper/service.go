package scraper

import (
	"context"
	"log"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

type srv struct {
}

type Product struct {
	ShopAdrr    string
	Category    string
	ProductUrl  string
	ProductName string
	SmallImg    string
	BigImg      string
	Price       string
	OldPrice    string
}

const (
	okeyMiaso = "https://www.okeydostavka.ru/spb/miaso-ptitsa-kolbasy/miaso-20"
	okeyURL   = "https://www.okeydostavka.ru/"
)

func NewService() *srv {
	return &srv{}
}

func (s *srv) Scrap() error {

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var good Product
	err := chromedp.Run(ctx, tasks(&good))

	if err != nil {
		return err
	}
	log.Printf("Go's time.After example:\n%s", good.ProductName)

	return nil
}

func tasks(g *Product) chromedp.Tasks {
	return chromedp.Tasks{
		network.Enable(),
		network.SetExtraHTTPHeaders(map[string]interface{}{"User-Agent": "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:126.0) Gecko/20100101 Firefox/126.0"}),
		chromedp.Navigate(okeyMiaso),
		chromedp.WaitVisible(`productListingWidget`),
		chromedp.Text(`.product-name`, &g.ProductName, chromedp.ByQueryAll),
	}
}

// func setHeaders(r *colly.Request) {
// 	r.Headers.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
// 	r.Headers.Set("Accept-Language", "en-GB,en;q=0.5")
// 	r.Headers.Set("Accept-Encoding", "gzip, deflate, br, zstd")
// 	r.Headers.Set("Connection", "keep-alive")
// 	r.Headers.Set("Host", "www.okeydostavka.ru")

// 	r.Headers.Set("DNT", "1")
// 	r.Headers.Set("Sec-GPC", "1")
// 	r.Headers.Set("Upgrade-Insecure-Requests", "1")
// 	r.Headers.Set("Sec-Fetch-Dest", "document")
// 	r.Headers.Set("Sec-Fetch-Mode", "navigate")
// 	r.Headers.Set("Sec-Fetch-Site", "none")
// 	r.Headers.Set("Sec-Fetch-User", "?1")
// 	r.Headers.Set("Priority", "u=1")
// }
