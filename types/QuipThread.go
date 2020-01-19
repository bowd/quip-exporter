package types

import (
	"github.com/kennygrant/sanitize"
	"strings"
)

type QuipThread struct {
	Thread          ThreadInfo
	HTML            string
	SharedFolderIDs []string
	UserIDs         []string
}

type ThreadInfo struct {
	ID          string
	Link        string
	Title       string
	CreatedUsec uint64
	UpdatedUsec uint64
	AuthorID    string
	Type        string
	Sharing     map[string]interface{}
}

func (qt QuipThread) Filename() string {
	return sanitize.Path(strings.Replace(qt.Thread.Title, "/", ":", -1))
}
