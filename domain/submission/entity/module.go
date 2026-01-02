package entity

type Module struct {
	ID      string
	Slug    string
	Title   string
	Grade   *Grade
	Subject *Subject
}

type Grade struct {
	ID   string
	Name string
}

type Subject struct {
	ID   string
	Name string
}
