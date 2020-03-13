package browser

import (
	"github.com/blevesearch/bleve"
	"github.com/bowd/quip-exporter/types"
)

type FolderResponse struct {
	*types.QuipFolder
	Path types.Path `json:"breadcrumbs"`
}

type ThreadResponse struct {
	*types.QuipThread
	Path types.Path `json:"breadcrumbs"`
}

type ThreadSearchResponse struct {
	Threads map[string]ThreadResponse `json:"threads"`
	Folders map[string]FolderResponse `json:"folders"`

	Meta *bleve.SearchResult `json:"meta"`
}

type CommentsResponse struct {
	Comments []*types.QuipMessage       `json:"comments"`
	Users    map[string]*types.QuipUser `json:"users"`
}
