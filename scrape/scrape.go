package scrape

import (
	"fmt"
	"log"
	"regexp"

	"github.com/gocolly/colly/v2"
)

type article struct {
	Title   string
	Summary string
	Author  string
	Date    string
	Images  []string
	Text    string
}

func newArticle(title, summary, author, date, text string, images []string) *article {
	return &article{
		Title:   title,
		Summary: summary,
		Author:  author,
		Date:    date,
		Images:  images,
		Text:    text,
	}
}

func ScrapeSmartBear(smartBearURL string) {
	var articles []article

	c := colly.NewCollector(
		colly.MaxDepth(2),
		colly.IgnoreRobotsTxt(),
		colly.UserAgent("Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36"),
	)

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Visited", r.Request.URL)
	})

	c.OnHTML("div.layout-container-content", func(h *colly.HTMLElement) {
		h.ForEach("ul", func(_ int, h *colly.HTMLElement) {
			h.ForEach("li", func(_ int, h *colly.HTMLElement) {
				h.ForEach("a[href]", func(_ int, h *colly.HTMLElement) {
					if h.Attr("href") != "/jason-cohen/" {
						h.Request.Visit(h.Attr("href"))
					}
				})
			})
		})
	})

	c.OnHTML("div.layout-container-content", func(h *colly.HTMLElement) {
		if h.Request.URL.Path != "/" {
			h.ForEach("article", func(_ int, h *colly.HTMLElement) {
				title := h.ChildText("h1")
				str := h.ChildText("div.print-biline")
				date, author, err := extractDateAndAuthor(str)
				if err != nil {
					log.Fatal(err)
				}
				var srcSets []string
				h.ForEach("div.e-content", func(_ int, h *colly.HTMLElement) {
					summary := h.ChildText("div.p-summary")
					text := h.Text
					h.ForEach("figure", func(_ int, h *colly.HTMLElement) {
						h.ForEach("div.img-container", func(_ int, h *colly.HTMLElement) {
							h.ForEach("picture", func(_ int, h *colly.HTMLElement) {
								srcSet := h.ChildAttr("source", "srcset")
								srcSets = append(srcSets, srcSet)
							})
						})
					})
					article := newArticle(title, summary, author, date, text, srcSets)
					articles = append(articles, *article)
				})
			})
		}
	})

	c.Visit(smartBearURL)

	fmt.Println("Articles:", len(articles))
}

func extractDateAndAuthor(str string) (string, string, error) {
	input := str

	pattern := `by (.+) on (\w+ \d{1,2}, \d{4})`

	re := regexp.MustCompile(pattern)

	matches := re.FindStringSubmatch(input)

	if len(matches) >= 3 {
		author := matches[1]
		date := matches[2]
		return date, author, nil
	}
	fmt.Println("No match found")
	return "", "", nil
}
