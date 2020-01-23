package scraper

import (
	"github.com/bowd/quip-exporter/client"
	"github.com/bowd/quip-exporter/types"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

type BlobNode struct {
	*BaseNode
	thread *types.QuipThread
	blob   []byte
}

func NewBlobNode(parent *ThreadNode, blobID string) INode {
	wg, ctx := errgroup.WithContext(parent.ctx)
	return &BlobNode{
		BaseNode: &BaseNode{
			logger: logrus.WithField("module", NodeTypes.Blob).
				WithField("id", blobID).
				WithField("path", parent.path),
			path: parent.path + "/blobs/",
			id:   blobID,
			ctx:  ctx,
			wg:   wg,
		},
		thread: parent.thread,
	}
}

func (node BlobNode) Type() NodeType {
	return NodeTypes.Blob
}

func (node BlobNode) ID() string {
	return node.id
}

func (node BlobNode) Children() []INode {
	return []INode{}
}

func (node *BlobNode) Process(scraper *Scraper) error {
	if node.ctx.Err() != nil {
		return nil
	}
	var blob []byte

	if exists, err := scraper.repo.BlobExists(node.thread.Thread.ID, node.id); err == nil && !exists {
		if blob, err := scraper.client.GetBlob(node.thread.Thread.ID, node.id); err != nil {
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
		if err := scraper.repo.SaveBlob(node.path, node.thread.Thread.ID, node.id, node.blob); err != nil {
			return err
		}
	} else if err != nil {
		node.logger.Errorln(err)
		return err
	} else {
		node.logger.Debugf("loaded from repository")
	}
	node.blob = blob
	return nil
}
