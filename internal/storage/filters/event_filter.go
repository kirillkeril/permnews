package filters

type EventFilter struct {
	Id       string
	Source   string
	Category string
	Skip     int
	Limit    int
	Title    string
}
