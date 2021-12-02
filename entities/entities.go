package entities

import "time"

type Todo struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreateAt  time.Time `json:"create_at"`
	IsChecked bool      `json:"is_checked"`
	Tags      []*Tag    `json:"tags"`
}

func (t *Todo) check() {
	t.IsChecked = true
}

type Tag struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
	Hex   string `json:"hex"`
}

type Checker interface {
	Check()
}
