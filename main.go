package main

import (
	"fmt"
	"io"
	"net/http"
	// "os"
)

type Feed struct {
    Url string
    Title string
    Description string
    Items []Item
}

type Item struct {
    Url string
    Title string
    Updated string
    Content string
}

func main() {
	db, err := makeDb();
	if err != nil {
		panic(err);
	}
	defer db.Close();

	err = createTables(db);
	if err != nil {
		panic(err);
	}

	for _,feedUrl := range getFeedsFromFile() {
		fmt.Printf("GETTING FEED: %v\n", feedUrl);
		feed, err := getGeneralFeedForm(feedUrl);
		if err != nil {
			panic(err);
		}
		fmt.Printf("FEED TITLE: %v\n", feed.Title);
		fmt.Printf("FEED URL: %v\n", feed.Url);
		fmt.Printf("FEED LEN ITEMS: %v\n", len(feed.Items));

		fmt.Printf("LEN FEEDS: %v\n", len(feed.Items));

		feedId, err := saveFeedToDb(db, feed);
		if err != nil {
			panic(err);
		}

		items, err := getItemsFromFeed(db, feed, feedId);
		if err != nil {
			panic(err);
		}
		fmt.Printf("DB LAST ITEM \n");
		fmt.Printf("DB TITLE %v\n", items[0].Title);
		fmt.Printf("DB URL %v\n", items[0].Url);
		
		fmt.Println("===================");
	}
}

func httpRequest(url string) (string, error) {
	resp, err := http.Get(url);
	if err != nil {
		return "", err;
	}
	defer resp.Body.Close();

	body, err := io.ReadAll(resp.Body);
	return string(body), nil;
}

func getFeedsFromFile() map[string]string {
	feeds := map[string]string{
		"j3s": "https://j3s.sh/feed.atom",
		"tadaima": "https://tadaima.bearblog.dev/feed",
		"ratfactor": "https://ratfactor.com/atom.xml",
	}
	return feeds;
}

