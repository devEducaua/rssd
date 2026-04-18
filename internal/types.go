package internal

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
