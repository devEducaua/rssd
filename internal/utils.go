package internal

import "time"

type Feed struct {
	Title string
	Name string
	Description string
	Url string
	Items []Item
}

type Item struct {
	Title string
	Updated string
	Content string
	Read bool
	Url string
}

func PeriodicReload(interval int) {
	req := []string{"UPDATE", "ALL"};
	updateCommand(req);

	ticker := time.NewTicker(time.Duration(interval));
	defer ticker.Stop();

	for range ticker.C {
		updateCommand(req);
	}
}
