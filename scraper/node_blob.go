package scraper

import (
	"github.com/bowd/quip-exporter/client"
	"github.com/bowd/quip-exporter/interfaces"
	"github.com/bowd/quip-exporter/types"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"path"
)

type BlobNode struct {
	*BaseNode
	thread *types.QuipThread
	blob   []byte
}

func NewBlobNode(parent *ThreadNode, blobID string) interfaces.INode {
	wg, ctx := errgroup.WithContext(parent.ctx)
	return &BlobNode{
		BaseNode: &BaseNode{
			logger: logrus.WithField("module", types.NodeTypes.Blob).
				WithField("id", blobID).
				WithField("path", parent.path),
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
	return []interfaces.INode{
		NewArchiveNode(
			path.Join(node.path, "blob", node.thread.Thread.ID),
			node.id,
			node.id,
			node,
		),
	}
}

func (node *BlobNode) Process(repo interfaces.IRepository, quip interfaces.IQuipClient) error {
	if node.ctx.Err() != nil {
		return nil
	}
	var blob []byte

	if exists, err := repo.NodeExists(node); err == nil && !exists {
		if blob, err := quip.GetBlob(node.thread.Thread.ID, node.id); err != nil {
			if client.IsUnauthorizedError(err) {
				node.logger.Warn("skipping unauthorised")
				return nil
			} else if err != nil {
				return err
			}
			return err
		} else {
			node.blob = blob
		}
		if err := repo.SaveNodeRaw(node, node.blob); err != nil {
			return err
		}
	} else if err != nil {
		node.logger.Errorln(err)
		return err
	} else {
		node.logger.Debugf("found in repo")
	}
	node.blob = blob
	return nil
}
