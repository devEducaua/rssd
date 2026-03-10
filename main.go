package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

var	feedsMap = make(map[string]string);

func main() {
	argv := os.Args[1:];

	if len(argv) == 0 {
		os.Exit(1);
	}

	switch argv[0] {
	case "update":
		fmt.Println("TODO: not implemented");
	case "list":
		listFeeds();
	case "add":
		addFeed(argv[1], argv[2]);
	default:
		getItemsByFeed(argv[0]);
	}
}

func addFeed(name string, url string) {
	feedsMap[name] = url;
}

func listFeeds() {
	fmt.Printf("FEEDS: %v\n", len(feedsMap));

	for name, url := range feedsMap {
		fmt.Printf("%v  ::  %v\n", name, url);
	}
}

func getItemsByFeed(feedName string) {
	url := feedsMap[feedName];

	feed, err := getGeneralFeedForm(url);
	if err != nil {
		fmt.Fprint(os.Stderr, err);
	}
	for _,v := range feed.Items {
		fmt.Printf("URL: %v\n", v.Title);
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
