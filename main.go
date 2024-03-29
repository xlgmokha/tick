package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

var (
	tickers map[string]string
)

func init() {
	tickers = map[string]string{}
}

func MapFrom(attrs []html.Attribute) map[string]string {
	m := map[string]string{}
	for _, attr := range attrs {
		m[attr.Key] = attr.Val
	}
	return m
}

func Visit(n *html.Node) {
	if n.Type == html.ElementNode && n.Data == "fin-streamer" {
		m := MapFrom(n.Attr)
		if m["data-field"] != "" && strings.Contains(m["data-field"], "Price") {
			tickers[m["data-symbol"]] = m["value"]
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		Visit(c)
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage:\n\t%s <symbol>\n", os.Args[0])
		return
	}
	ticker := strings.ToUpper(os.Args[1])
	response, _ := http.Get(fmt.Sprintf("https://finance.yahoo.com/quote/%s/", ticker))
	doc, _ := html.Parse(response.Body)
	Visit(doc)
	fmt.Printf("%v\n", tickers[ticker])
}
