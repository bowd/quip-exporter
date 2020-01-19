package interfaces

import "github.com/bowd/quip-exporter/types"

type IQuipClient interface {
	GetFolder(folderID string) (*types.QuipFolder, error)
	GetThread(threadID string) (*types.QuipThread, error)
}
