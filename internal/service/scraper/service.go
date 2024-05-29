package scraper

import (
	"context"
	"fmt"
	"log"
	"path"
	"time"

	"github.com/chromedp/cdproto/cdp"
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
	userAgent = "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:126.0) Gecko/20100101 Firefox/126.0"
	okeyMiaso = "https://www.okeydostavka.ru/spb/miaso-ptitsa-kolbasy/miaso-20"

	okeyURL    = "https://www.okeydostavka.ru/"
	okeyImgUrl = "https://36ltwco2hg.a.trbcdn.net"
)

func NewService() *srv {
	return &srv{}
}

func (s *srv) Scrap() error {
	ctx, _ := chromedp.NewContext(context.Background())
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := chromedp.Run(ctx, tasks()); err != nil {
		return fmt.Errorf("could not navigate %s", err.Error())
	}

	if err := chromedp.Run(ctx, chromedp.Sleep(2000*time.Millisecond) /*chromedp.WaitVisible("productListingWidget")*/); err != nil {
		return fmt.Errorf("could not get section %s", err.Error())
	}

	var productsNode []*cdp.Node
	if err := chromedp.Run(ctx, chromedp.Nodes(`.ok-theme.product`, &productsNode, chromedp.ByQueryAll)); err != nil {
		return fmt.Errorf("could not get nodes %s", err.Error())
	}

	goods := make([]Product, 0, len(productsNode))

	for _, v := range productsNode {
		var product Product
		//Get name
		chromedp.Run(ctx, chromedp.Text(".product-name", &product.ProductName, chromedp.ByQuery, chromedp.FromNode(v)))

		var ok bool

		//Get product url
		chromedp.Run(ctx, chromedp.AttributeValue(".product-name > a", "href", &product.ProductUrl, &ok, chromedp.ByQuery, chromedp.FromNode(v)))
		if !ok {
			log.Print("could not get product url")
		}
		product.ProductUrl = path.Join(okeyURL, product.ProductUrl)

		//Get small image url
		chromedp.Run(ctx, chromedp.AttributeValue(".product-image img", "data-src", &product.SmallImg, &ok, chromedp.ByQuery, chromedp.FromNode(v)))
		if !ok {
			log.Print("could not get small image url")
		}
		product.SmallImg = path.Join(okeyImgUrl, product.SmallImg)
		product.BigImg = product.SmallImg

		//Get old price
		chromedp.Run(ctx, chromedp.Text(".label.crossed", &product.OldPrice, chromedp.ByQuery, chromedp.FromNode(v)))

		//Get current price
		chromedp.Run(ctx, chromedp.Text(".label.price:last-child", &product.Price, chromedp.ByQuery, chromedp.FromNode(v)))

		goods = append(goods, product)
	}

	for _, v := range goods {
		fmt.Println(v)
	}
	return nil
}

func tasks() chromedp.Tasks {
	return chromedp.Tasks{
		network.Enable(),
		network.SetExtraHTTPHeaders(map[string]interface{}{"User-Agent": userAgent}),
		chromedp.Navigate(okeyMiaso),
	}
}
