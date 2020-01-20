package interfaces

import "github.com/bowd/quip-exporter/types"

type IRepository interface {
	GetCurrentUser() (*types.QuipUser, error)
	GetUser(id string) (*types.QuipUser, error)
	GetFolder(id string) (*types.QuipFolder, error)
	GetThread(id string) (*types.QuipThread, error)
}
