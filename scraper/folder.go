package scraper

import (
	"context"
	"fmt"
	"github.com/bowd/quip-exporter/client"
	"github.com/bowd/quip-exporter/repositories"
	"github.com/bowd/quip-exporter/types"
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

func (node *FolderNode) ID() string {
	return fmt.Sprintf("folder:%s [%s]", node.id, node.path)
}

func (node *FolderNode) Children() []INode {
	if node.folder == nil {
		return []INode{}
	}

	children := make([]INode, 0, 0)
	nodePath := path.Join(node.path, node.folder.PathSegment())
	for _, child := range node.folder.Children {
		if child.IsThread() {
			childNode := NewThreadNode(node.ctx, nodePath, *child.ThreadID)
			children = append(children, childNode)
		} else if child.IsFolder() {
			childNode := NewFolderNode(node.ctx, nodePath, *child.FolderID)
			children = append(children, childNode)
		}
	}
	return children
}

func (node *FolderNode) Process(scraper *Scraper) error {
	if node.ctx.Err() != nil {
		return nil
	}
	var folder *types.QuipFolder
	var err error

	folder, err = scraper.repo.GetFolder(node.id)
	if err != nil && repositories.IsNotFoundError(err) {
		folder, err = scraper.client.GetFolder(node.id)
		if err != nil && client.IsUnauthorizedError(err) {
			scraper.logger.Debugf("Ignoring Unauthorized repo:%s [%s]", node.id, node.path)
			return nil
		} else if err != nil {
			return err
		}
		if err := scraper.repo.SaveFolder(folder); err != nil {
			return err
		}
	} else if err != nil {
		return err
	} else {
		scraper.logger.Debugf("Loaded from repo folder:%s [%s]", node.id, node.path)
	}

	node.folder = folder
	return nil
}
