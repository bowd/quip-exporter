package path

import (
	"context"
	"github.com/bowd/quip-exporter/interfaces"
	"github.com/bowd/quip-exporter/scraper"
	"github.com/bowd/quip-exporter/types"
)

func getParentStack(c context.Context, parentID string, repo interfaces.IRepository) (types.Path, error) {
	parent, err := repo.GetFolder(scraper.NewFolderNode(c, "/", parentID))
	if err != nil {
		return nil, err
	}

	parentPathItem := types.PathFolder{
		FolderID: parentID,
		Title:    parent.Folder.Title,
	}

	stack, err := GetPathToFolder(c, parent, repo)
	if err != nil {
		return nil, err
	}
	return append(stack, parentPathItem), nil
}

func GetPathToFolder(c context.Context, folder *types.QuipFolder, repo interfaces.IRepository) (types.Path, error) {
	if folder.Folder.ParentID == "" {
		return types.Path{}, nil
	}
	return getParentStack(c, folder.Folder.ParentID, repo)

}

func GetPathToThread(c context.Context, thread *types.QuipThread, repo interfaces.IRepository) (types.Path, error) {
	if len(thread.SharedFolderIDs) == 0 {
		return types.Path{}, nil
	}
	return getParentStack(c, thread.SharedFolderIDs[0], repo)
}
