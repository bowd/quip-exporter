package types

type QuipMessage struct {
	ID          string                `json:"id"`
	AuthorID    string                `json:"author_id"`
	CreatedUsec uint64                `json:"created_usec"`
	Text        string                `json:"text"`
	Parts       []interface{}         `json:"parts"`
	Annotation  QuipMessageAnnotation `json:"annotation"`
}

type QuipMessageAnnotation struct {
	ID                  string   `json:"id"`
	HighlightSectionIDs []string `json:"highlight_section_ids"`
}
