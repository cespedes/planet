package feed

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/mmcdole/gofeed"
	"golang.org/x/net/html"
)

// TODO: do not fetch URL twice...
func GetFeed(addr string) (string, error) {
	resp, err := http.Get(addr)
	if err != nil {
		return "", err
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return "", err
	}

	type link struct {
		Rel  string
		Type string
		Href string
	}
	var links []link
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "link" {
			var l link
			for _, a := range n.Attr {
				switch a.Key {
				case "rel":
					l.Rel = a.Val
				case "type":
					l.Type = a.Val
				case "href":
					l.Href = a.Val
				}
			}
			links = append(links, l)
			//log.Printf("<link attrs=%v>\n", n.Attr)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	fmt.Printf("Links: %+v\n", links)
	for _, l := range links {
		if (l.Rel == "alternate" && l.Type == "application/atom+xml") ||
			(l.Rel == "alternate" && l.Type == "application/rss+xml") ||
			(l.Rel == "alternate" && l.Type == "text/xml+oembed") {
			u, err := url.Parse(l.Href)
			if err != nil {
				log.Fatal(err)
			}
			base, err := url.Parse(addr)
			if err != nil {
				log.Fatal(err)
			}
			//log.Printf("GetFeed: base=%q u=%q", base, u)
			href := base.ResolveReference(u).String()
			//log.Printf("GetFeed: returning %q", href)
			return href, nil
		}
	}
	return "", fmt.Errorf("no feed found")
}

// Parse parses a URL and returns a feed if it is found.
func Parse(addr string) (*gofeed.Feed, error) {
	parser := gofeed.NewParser()
	feed, err := parser.ParseURL(addr)
	if err == nil {
		return feed, nil
	}
	alternate, err := GetFeed(addr)
	if err == nil {
		feed, err = parser.ParseURL(alternate)
		if err == nil {
			return feed, err
		}
	}

	//return nil, fmt.Errorf("XX1: %s: %v", alternate, err)

	base, err := url.Parse(addr)
	if err != nil {
		log.Fatal(err)
	}
	u, err := url.Parse("index.xml")
	if err != nil {
		log.Fatal(err)
	}
	//log.Printf("GetFeed: base=%q u=%q", base, u)
	feed, err = parser.ParseURL(base.ResolveReference(u).String())

	if err != nil {
		if alternate != "" {
			addr = alternate
		}
		return nil, fmt.Errorf("%s: %v", addr, err)
	}
	return feed, nil
}
