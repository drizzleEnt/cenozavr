package service

type Scraper interface {
	Scrap() error
}
