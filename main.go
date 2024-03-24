package main

import "github.com/haron1996/tiktok-scrapper/scrape"

func main() {
	smartBearURL := "https://longform.asmartbear.com/"
	scrape.ScrapeSmartBear(smartBearURL)
}
