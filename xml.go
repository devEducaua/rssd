package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type XmlAtomFeed struct {
    Id string `xml:"id"`
    Title string `xml:"title"`
    Subtitle string `xml:"subtitle"`
    Entries []XmlAtomEntry `xml:"entry"`
}

type XmlAtomEntry struct {
    Id string `xml:"id"`
    Title string `xml:"title"`
    Updated string `xml:"updated"`
    Content string `xml:"content"`
}

type XmlRss struct {
	Channel XmlRssFeed `xml:"channel"`
}

type XmlRssFeed struct {
    Id string `xml:"id"`
    Title string `xml:"title"`
    Description string `xml:"description"`
    Items []XmlRssItem `xml:"item"`
}

type XmlRssItem struct {
    Id string `xml:"id"`
    Title string `xml:"title"`
    PubDate string `xml:"pubDate"`
    Description string `xml:"description"`
}

func getFeedFromWeb(feedUrl string) (Feed, error) {
	u, err := url.Parse(feedUrl);
	if err != nil {
		return Feed{}, err;
	}

	var rawXml string;
	switch u.Scheme {
	case "gemini":
		rawXml, err = geminiRequest(feedUrl);
	case "https", "http":
		rawXml, err = httpRequest(feedUrl);
	default:
		return Feed{}, fmt.Errorf("not supported scheme");
	}

	// TODO: support RSS
	feed, err := atomToGenericForm(rawXml);
	if err != nil {
		return Feed{}, nil;
	}

	return feed, nil;
}

func httpRequest(url string) (string, error) {
	res, err := http.Get(url);
	if err != nil {
		return "", err;
	}
	defer res.Body.Close();

	body, err := io.ReadAll(res.Body);
	if err != nil {
		return "", err;
	}

	return string(body), nil;
}

func geminiRequest(url string) (string, error) {
	panic("TODO: implement gemini requests");
}

func atomToGenericForm(xmlFile string) (Feed, error) {
    var atom XmlAtomFeed;

    err := xml.Unmarshal([]byte(xmlFile), &atom);
    if err != nil {
        return Feed{}, fmt.Errorf("ERROR: could not parse the atom file: %v\n", err);
    }
    var items []Item

    for _, e := range atom.Entries {
        items = append(items, Item{
            Url:        e.Id,
            Title:     e.Title,
            Updated: e.Updated,
            Content:   e.Content,
			Read: false,
        })
    }
    feed := Feed{
        Url: atom.Id,
        Title: atom.Title,
		Name: atom.Title,
        Description: atom.Subtitle,
        Items: items,
    }

    return feed, nil;
}
