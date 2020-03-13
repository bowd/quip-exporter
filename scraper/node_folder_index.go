package scraper

import (
	"github.com/bowd/quip-exporter/interfaces"
	"github.com/bowd/quip-exporter/types"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"path"
)

type FolderIndexNode struct {
	*FolderNode
	exists bool
}

func NewFolderIndexNode(parent *FolderNode) interfaces.INode {
	wg, ctx := errgroup.WithContext(parent.ctx)
	return &FolderIndexNode{
		FolderNode: &FolderNode{
			BaseNode: &BaseNode{
				logger: logrus.WithField("module", types.NodeTypes.FolderIndex).
					WithField("id", parent.id).
					WithField("path", parent.path),
				path: parent.path,
				id:   parent.id,
				ctx:  ctx,
				wg:   wg,
			},
			folder: parent.folder,
		},
	}
}

func (node FolderIndexNode) Type() types.NodeType {
	return types.NodeTypes.FolderIndex
}

func (node *FolderIndexNode) ID() string {
	return node.id + "/index"
}

func (node FolderIndexNode) Path() string {
	return path.Join("data", "index", node.id)
}

func (node *FolderIndexNode) Children() []interfaces.INode {
	return []interfaces.INode{}
}

func (node *FolderIndexNode) Process(repo interfaces.IRepository, quip interfaces.IQuipClient, search interfaces.ISearchIndex) error {
	if node.ctx.Err() != nil {
		return nil
	}

	if indexed, err := search.IsIndexed(node.folder.Folder.ID); err != nil {
		return err
	} else if indexed {
		node.logger.Debug("Skipping already indexed")
		return nil
	}

	return search.IndexFolder(
		node.folder.Folder.ID,
		node.folder.Folder.Title,
	)
}
