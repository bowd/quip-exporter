package interfaces

import "github.com/bowd/quip-exporter/types"

type IQuipClient interface {
	GetFolder(folderID string) (*types.QuipFolder, error)
	GetThread(threadID string) (*types.QuipThread, error)
	GetCurrentUser() (*types.QuipUser, error)
	GetUser(userID string) (*types.QuipUser, error)

	ExportThreadSlides(threadID string) ([]byte, error)
	ExportThreadDocument(threadID string) ([]byte, error)
	ExportThreadSpreadsheet(threadID string) ([]byte, error)
	ExportUserPhoto(url string) ([]byte, error)

	GetThreadComments(threadID string) ([]*types.QuipMessage, error)
	GetBlob(threadID, blobID string) ([]byte, error)
}
