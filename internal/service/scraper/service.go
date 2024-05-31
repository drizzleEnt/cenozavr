package scraper

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"
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
		err := chromedp.Run(ctx, chromedp.Text(".product-name", &product.ProductName, chromedp.ByQuery, chromedp.FromNode(v)))
		if err != nil {
			log.Print("could not get product name")
		}

		var ok bool

		//Get product url
		err = chromedp.Run(ctx, chromedp.AttributeValue(".product-name > a", "href", &product.ProductUrl, &ok, chromedp.ByQuery, chromedp.FromNode(v)))
		if err != nil {
			log.Print("could not get product name")
		}
		product.ProductUrl = path.Join(okeyURL, product.ProductUrl)

		//Get small image url
		err = chromedp.Run(ctx, chromedp.AttributeValue(".product-image img", "data-src", &product.SmallImg, &ok, chromedp.ByQuery, chromedp.FromNode(v)))
		if err != nil {
			log.Print("could not get product image url")
		}
		product.SmallImg = path.Join(okeyImgUrl, product.SmallImg)
		product.BigImg = product.SmallImg

		//Get old price
		err = chromedp.Run(ctx, chromedp.Text(".label.crossed", &product.OldPrice, chromedp.ByQuery, chromedp.FromNode(v)))
		if err != nil {
			log.Print("could not get product price")
		}
		//Get current price
		err = chromedp.Run(ctx, chromedp.Text(".label.price:last-child", &product.Price, chromedp.ByQuery, chromedp.FromNode(v)))
		if err != nil {
			log.Print("could not get product current price")
		}

		goods = append(goods, product)
	}

	if err := save(goods); err != nil {
		return err
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

func save(products []Product) error {
	file, err := os.Create("./bin/products.csv")
	if err != nil {
		return err
	}

	defer file.Close()

	w := csv.NewWriter(file)
	defer w.Flush()

	headers := []string{
		"name",
		"url",
		"small imgUrl",
		"big imgUrl",
		"current price",
		"old price",
	}

	err = w.Write(headers)
	if err != nil {
		return err
	}

	for _, p := range products {
		record := []string{
			p.ProductName,
			p.ProductUrl,
			p.SmallImg,
			p.BigImg,
			p.Price,
			p.OldPrice,
		}
		err := w.Write(record)
		if err != nil {
			return err
		}
	}

	return nil
}
