package scraper

import (
	"context"
	"github.com/bowd/quip-exporter/client"
	"github.com/bowd/quip-exporter/repositories"
	"github.com/bowd/quip-exporter/types"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"path"
)

type FolderNode struct {
	*BaseNode
	folder *types.QuipFolder
}

func NewFolderNode(ctx context.Context, path, id string) INode {
	wg, ctx := errgroup.WithContext(ctx)
	return &FolderNode{
		BaseNode: &BaseNode{
			logger: logrus.WithField("module", NodeTypes.Folder).
				WithField("id", id).
				WithField("path", path),
			id:   id,
			path: path,
			wg:   wg,
			ctx:  ctx,
		},
	}
}

func (node FolderNode) Type() NodeType {
	return NodeTypes.Folder
}

func (node FolderNode) Children() []INode {
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
			node.logger.Warn("skipping unauthorised")
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
		node.logger.Debugf("loaded from repository")
	}

	node.folder = folder
	return nil
}
