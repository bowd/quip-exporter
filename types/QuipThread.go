package types

import (
	"github.com/kennygrant/sanitize"
	"strings"
)

type QuipThread struct {
	Thread          ThreadInfo `json:"thread"`
	HTML            string     `json:"html"`
	SharedFolderIDs []string   `json:"shared_folder_ids"`
	UserIDs         []string   `json:"user_ids"`
}

type ThreadInfo struct {
	ID          string                 `json:"id"`
	Link        string                 `json:"link"`
	Title       string                 `json:"title"`
	CreatedUsec uint64                 `json:"created_usec"`
	UpdatedUsec uint64                 `json:"updated_usec"`
	AuthorID    string                 `json:"author_id"`
	Type        string                 `json:"type"`
	Sharing     map[string]interface{} `json:"sharing"`
}

func (qt QuipThread) Filename() string {
	return sanitize.Path(strings.Replace(qt.Thread.Title, "/", ":", -1))
}

func (qt QuipThread) IsSpreadsheet() bool {
	return qt.Thread.Type == "spreadsheet"
}

func (qt QuipThread) IsSlides() bool {
	return qt.Thread.Type == "slides"
}

func (qt QuipThread) IsChannel() bool {
	return qt.Thread.Type == "channel"
}

func (qt QuipThread) IsDocument() bool {
	return qt.Thread.Type == "document"
}
