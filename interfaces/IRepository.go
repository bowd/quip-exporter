package interfaces

import "github.com/bowd/quip-exporter/types"

type IRepository interface {
	GetCurrentUser() (*types.QuipUser, error)
	GetUser(id string) (*types.QuipUser, error)
	GetFolder(id string) (*types.QuipFolder, error)
	GetThread(id string) (*types.QuipThread, error)
	SaveCurrentUser(user *types.QuipUser) error
	SaveUser(user *types.QuipUser) error
	SaveThread(thread *types.QuipThread) error
	SaveFolder(folder *types.QuipFolder) error

	HasExportedHTML(threadID string) (bool, error)
	SaveThreadHTML(nodePath string, thread *types.QuipThread) error

	HasExportedSlides(threadID string) (bool, error)
	SaveThreadSlides(nodePath string, thread *types.QuipThread, pdf []byte) error

	HasExportedDocument(threadID string) (bool, error)
	SaveThreadDocument(nodePath string, thread *types.QuipThread, doc []byte) error

	HasExportedSpreadsheet(threadID string) (bool, error)
	SaveThreadSpreadsheet(nodePath string, thread *types.QuipThread, xls []byte) error

	GetThreadComments(threadID string) ([]*types.QuipMessage, error)
	SaveThreadComments(threadID string, comments []*types.QuipMessage) error
}
