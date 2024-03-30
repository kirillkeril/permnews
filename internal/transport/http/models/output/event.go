package output

type Event struct {
	Id       string `json:"id"`
	Title    string `json:"title"`
	Notice   string `json:"notice"`
	Source   string `json:"source"`
	Image    string `json:"image"`
	Link     string `json:"link"`
	Category string `json:"category"`
}
