package browser

import "github.com/bowd/quip-exporter/types"

type FolderResponse struct {
	*types.QuipFolder
	Path types.Path `json:"breadcrumbs"`
}

type ThreadResponse struct {
	*types.QuipThread
	Path types.Path `json:"breadcrumbs"`
}
