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

func rssToGeneral(xmlFile string) Feed {
    var rss XmlRss;

    err := xml.Unmarshal([]byte(xmlFile), &rss);
    if err != nil {
        fmt.Fprintf(os.Stderr, "ERROR: could not parse the rss file: %v\n", err);
        os.Exit(1);
    }

    var items []Item;
    for _, e := range rss.Channel.Items {
        items = append(items, Item{
            Url:        e.Id,
            Title:     e.Title,
            Updated: e.PubDate,
            Content:   e.Description,
        })
    }

    feed := Feed{
        Url: rss.Channel.Id,
        Title: rss.Channel.Title,
        Description: rss.Channel.Description,
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
            Url:        e.Id,
            Title:     e.Title,
            Updated: e.Updated,
            Content:   e.Content,
        })
    }
    feed := Feed{
        Url: atom.Id,
        Title: atom.Title,
        Description: atom.Subtitle,
        Items: items,
    }

    return feed, nil;
}

func requestFeed(feedUrl string) (Feed, error) {
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
