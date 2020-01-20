package scraper

import (
	"context"
	"fmt"
	"github.com/bowd/quip-exporter/types"
	"github.com/bowd/quip-exporter/utils"
	"golang.org/x/sync/errgroup"
	"path"
)

type FolderNode struct {
	*Node
	folder *types.QuipFolder
}

func NewFolderNode(ctx context.Context, path, id string) INode {
	wg, ctx := errgroup.WithContext(ctx)
	return &FolderNode{
		Node: &Node{
			id:   id,
			path: path,
			wg:   wg,
			ctx:  ctx,
		},
	}
}

func (node *FolderNode) Children() []INode {
	if node.folder == nil {
		return []INode{}
	}

	children := make([]INode, 0, 0)
	nodePath := path.Join(node.path, node.folder.Folder.Title)
	for _, child := range node.folder.Children {
		if child.IsThread() {
			children = append(children, NewThreadNode(node.ctx, nodePath, *child.ThreadID))
		} else if child.IsFolder() {
			children = append(children, NewFolderNode(node.ctx, nodePath, *child.FolderID))
		}
	}
	return children
}

func (node *FolderNode) Load(scraper *Scraper) error {
	scraper.logger.Debugf("Handling folder:%s [%s]", node.id, node.path)
	if node.ctx.Err() != nil {
		return nil
	}

	folder, err := scraper.client.GetFolder(node.id)
	if err != nil {
		return fmt.Errorf("could not get folder: %s: %s", node.id, err)
	}
	node.folder = folder

	if err := node.saveData(scraper); err != nil {
		return err
	}
	return nil
}

func (node *FolderNode) saveData(scraper *Scraper) error {
	if node.ctx.Err() != nil {
		return nil
	}
	if err := utils.SaveJSONToFile(path.Join(DATA_FOLDER, "folders", node.id+".json"), node.folder); err != nil {
		return err
	}
	return nil
}
