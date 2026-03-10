package main

import (
    "encoding/xml"
    "fmt"
    "os"
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
    Published string `xml:"published"`
    Content string `xml:"content"`
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

type Feed struct {
    Id string
    Title string
    Description string
    Items []Item
}

type Item struct {
    Id string
    Title string
    Published string
    Content string
}

func rssToGeneral(xmlFile string) Feed {
    var rss XmlRssFeed;

    err := xml.Unmarshal([]byte(xmlFile), &rss);
    if err != nil {
        fmt.Fprintf(os.Stderr, "ERROR: could not parse the rss file: %v\n", err);
        os.Exit(1);
    }

    var items []Item;
    for _, e := range rss.Items {
        items = append(items, Item{
            Id:        e.Id,
            Title:     e.Title,
            Published: e.PubDate,
            Content:   e.Description,
        })
    }

    feed := Feed{
        Id: rss.Id,
        Title: rss.Title,
        Description: rss.Description,
        Items: items,
    }

    return feed;
}

func atomToGeneral(xmlFile string) (Feed, error) {
    var atom XmlAtomFeed;

    err := xml.Unmarshal([]byte(xmlFile), &atom);
    if err != nil {
        return Feed{}, fmt.Errorf("ERROR: could not parse the atom file: %v\n", err);
    }
    var items []Item

    for _, e := range atom.Entries {
        items = append(items, Item{
            Id:        e.Id,
            Title:     e.Title,
            Published: e.Published,
            Content:   e.Content,
        })
    }
    feed := Feed{
        Id: atom.Id,
        Title: atom.Title,
        Description: atom.Subtitle,
        Items: items,
    }

    return feed, nil;
}

func getGeneralFeedForm(feedUrl string) (Feed, error) {
	// TODO: handle more protocols like: gemini and gopher.
    xmlFile, err := httpRequest(feedUrl);
    if err != nil {
        return Feed{}, fmt.Errorf("ERROR: could not request the file: %v\n", err);
    }

	// TODO: this handles only atom;
    var feed Feed;
	feed, err = atomToGeneral(xmlFile);
	if err != nil {
		return Feed{}, err;
	}

    return feed, nil;
}
