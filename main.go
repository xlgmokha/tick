package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

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
			fmt.Printf("%v: %s\n", m["data-symbol"], m["value"])
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		Visit(c)
	}
}

func main() {
	response, err1 := http.Get("https://finance.yahoo.com/quote/MSFT/")
	if err1 != nil {
		os.Exit(1)
	}

	doc, err := html.Parse(response.Body)
	if err != nil {
		os.Exit(2)
	}
	Visit(doc)
}
