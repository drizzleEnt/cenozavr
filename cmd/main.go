package main

import (
	"log"

	"github.com/drizzleent/cenozavr/internal/service"
	"github.com/drizzleent/cenozavr/internal/service/scraper"
)

const (
	userAgent = "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:125.0) Gecko/20100101 Firefox/125.0"
)

var s service.Scraper

func main() {

	s = scraper.NewService()

	err := s.Scrap()
	if err != nil {
		log.Printf("[ERROR]: %s", err.Error())
	}

}
