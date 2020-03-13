package scraper

import (
	"github.com/bowd/quip-exporter/interfaces"
	"github.com/bowd/quip-exporter/types"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"path"
)

type BlobNode struct {
	*BaseNode
	thread *types.QuipThread
	exists bool
}

func NewBlobNode(parent *ThreadNode, blobID string) interfaces.INode {
	wg, ctx := errgroup.WithContext(parent.ctx)
	return &BlobNode{
		BaseNode: &BaseNode{
			logger: logrus.WithField("module", types.NodeTypes.Blob).
				WithField("id", blobID).
				WithField("path", parent.path).
				WithField("thread", parent.id),
			path: path.Join(parent.path, "blob"),
			id:   blobID,
			ctx:  ctx,
			wg:   wg,
		},
		thread: parent.thread,
	}
}

func (node BlobNode) Type() types.NodeType {
	return types.NodeTypes.Blob
}

func (node BlobNode) ID() string {
	return node.thread.Thread.ID + "/" + node.id
}

func (node BlobNode) Path() string {
	return path.Join("data", "blobs", node.thread.Thread.ID, node.id)
}

func (node *BlobNode) Children() []interfaces.INode {
	if !node.exists {
		return []interfaces.INode{}
	}
	return []interfaces.INode{
		NewArchiveNode(
			path.Join(node.path, node.thread.Thread.ID),
			node.id,
			node.id,
			node,
		),
	}
}

func (node *BlobNode) Process(repo interfaces.IRepository, quip interfaces.IQuipClient, search interfaces.ISearchIndex) error {
	if node.ctx.Err() != nil {
		return nil
	}

	if exists, err := repo.NodeExists(node); err == nil && !exists {
		if blob, err := quip.GetBlob(node.thread.Thread.ID, node.id); err != nil {
			return err
		} else {
			if err := repo.SaveNodeRaw(node, blob); err != nil {
				return err
			} else {
				node.exists = true
			}
		}
	} else if err != nil {
		return err
	} else {
		node.exists = true
		node.logger.Debugf("found in repo")
	}
	return nil
}
