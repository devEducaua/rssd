package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
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
	Read bool
}

const FEEDSFILEPATH = "./feeds.txt";	
const SOCKPATH = "/tmp/rssd.sock";

func main() {
	if err := os.Remove(SOCKPATH); err != nil && !os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "ERROR: failed to remove the file %v: %v\n", SOCKPATH, err);
		os.Exit(1);
	}

	listener, err := net.Listen("unix", SOCKPATH);
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: on listen: %v\n", err);
		os.Exit(1);
	}

	defer listener.Close();

	for {
		conn, err := listener.Accept();
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: on accept: %v\n", err);
		}

		go handleConnection(conn);
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close();
	fmt.Printf("REMOTE: %v\n", conn.LocalAddr().String());
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

func parseFeedsFile(feedsFilePath string) (map[string]string, error) {
	feeds := make(map[string]string);

	bytes, err := os.ReadFile(FEEDSFILEPATH);
	if err != nil {
		return nil, err;
	}

	lines := strings.Split(string(bytes), "\n");	

	for _,l := range lines {
		if strings.TrimSpace(l) != "" {
			parts := strings.SplitN(l, " ", 2);
			feeds[parts[0]] = parts[1];
		}
	}

	return feeds, nil;
}

