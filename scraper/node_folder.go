package scraper

import (
	"context"
	"github.com/bowd/quip-exporter/interfaces"
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

func NewFolderNode(ctx context.Context, path, id string) interfaces.INode {
	wg, ctx := errgroup.WithContext(ctx)
	return &FolderNode{
		BaseNode: &BaseNode{
			logger: logrus.WithField("module", types.NodeTypes.Folder).
				WithField("id", id).
				WithField("path", path),
			id:   id,
			path: path,
			wg:   wg,
			ctx:  ctx,
		},
	}
}

func (node FolderNode) Path() string {
	return path.Join("data", "folders", node.id+".json")
}

func (node FolderNode) Type() types.NodeType {
	return types.NodeTypes.Folder
}

func (node FolderNode) Children() []interfaces.INode {
	if node.folder == nil {
		return []interfaces.INode{}
	}

	children := make([]interfaces.INode, 0, 0)
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

func (node *FolderNode) Process(repo interfaces.IRepository, quip interfaces.IQuipClient, search interfaces.ISearchIndex) error {
	if node.ctx.Err() != nil {
		return nil
	}
	var folder *types.QuipFolder
	var err error

	folder, err = repo.GetFolder(node)
	if err != nil && repositories.IsNotFoundError(err) {
		if folder, err = quip.GetFolder(node.id); err != nil {
			return err
		}
		if err := repo.SaveNodeJSON(node, folder); err != nil {
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
