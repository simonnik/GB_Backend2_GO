package links

import (
	"errors"
	"net/url"
	"time"
)

type Link struct {
	ID    int64  `json:"-"`
	Link  string `json:"link"`
	Token string `json:"-"`
}

type Stat struct {
	ID      int64  `json:"id"`
	Link    string `json:"link"`
	IP      string `json:"ip"`
	Created string `json:"created_at" db:"created_at"`
}

func (s *Stat) FormattedDate() string {
	t, _ := time.Parse(time.RFC3339, s.Created)

	return t.Format(time.RFC822)
}

type StatList []*Stat

func (l Link) Validate() error {
	if l.Link == "" {
		return errors.New("link can't be empty")
	}

	if _, err := url.ParseRequestURI(l.Link); err != nil {
		return errors.New("link is invalid")
	}

	return nil
}
