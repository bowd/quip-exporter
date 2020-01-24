package interfaces

import "github.com/bowd/quip-exporter/types"

type IRepository interface {
	GetCurrentUser(INode) (*types.QuipUser, error)
	GetUser(INode) (*types.QuipUser, error)
	GetFolder(INode) (*types.QuipFolder, error)
	GetThread(INode) (*types.QuipThread, error)
	GetThreadComments(INode) ([]*types.QuipMessage, error)
	NodeExists(INode) (bool, error)
	SaveNodeJSON(INode, interface{}) error
	SaveNodeRaw(INode, []byte) error
	MakeArchiveCopy(INode, INode) error
}
