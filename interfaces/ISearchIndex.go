package interfaces

type ISearchIndex interface {
	IndexThread(id, title, content string) error
	IndexFolder(id, title string) error
	IsIndexed(id string) (bool, error)
}
