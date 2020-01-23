package scraper

import (
	"fmt"
	"github.com/bowd/quip-exporter/types"
	"golang.org/x/sync/errgroup"
)

type BlobNode struct {
	*Node
	thread *types.QuipThread
	blob   []byte
}

func NewBlobNode(parent *ThreadNode, blobID string) INode {
	wg, ctx := errgroup.WithContext(parent.ctx)
	return &BlobNode{
		Node: &Node{
			path: parent.path + "/blobs/",
			id:   blobID,
			ctx:  ctx,
			wg:   wg,
		},
		thread: parent.thread,
	}
}

func (node BlobNode) ID() string {
	return fmt.Sprintf("blob:%s [%s]", node.id, node.path)
}

func (node BlobNode) Children() []INode {
	return []INode{}
}

func (node *BlobNode) Process(scraper *Scraper) error {
	scraper.logger.Debugf("Handling blob:%s", node.id)
	if node.ctx.Err() != nil {
		return nil
	}
	var blob []byte

	if exists, err := scraper.repo.BlobExists(node.thread.Thread.ID, node.id); err == nil && !exists {
		if blob, err := scraper.client.GetBlob(node.thread.Thread.ID, node.id); err != nil {
			return err
		} else {
			node.blob = blob
		}
		if err := scraper.repo.SaveBlob(node.path, node.thread.Thread.ID, node.id, node.blob); err != nil {
			return err
		}
	} else if err != nil {
		return err
	} else {
		scraper.logger.Debugf("Loaded from repo blob:%s", node.id)
	}
	node.blob = blob
	return nil
}
