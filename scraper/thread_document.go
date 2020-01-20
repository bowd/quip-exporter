package scraper

import (
	"fmt"
	"golang.org/x/sync/errgroup"
)

type ThreadDocumentNode struct {
	*ThreadNode
}

func NewThreadDocumentNode(parent *ThreadNode) INode {
	wg, ctx := errgroup.WithContext(parent.ctx)
	return &ThreadDocumentNode{
		ThreadNode: &ThreadNode{
			Node: &Node{
				path: parent.path,
				id:   parent.id,
				ctx:  ctx,
				wg:   wg,
			},
			thread: parent.thread,
		},
	}
}

func (node *ThreadDocumentNode) ID() string {
	return fmt.Sprintf("thread:%s:document [%s]", node.id, node.path)
}

func (node *ThreadDocumentNode) Children() []INode {
	return []INode{}
}

func (node *ThreadDocumentNode) Process(scraper *Scraper) error {
	if node.ctx.Err() != nil {
		return nil
	}
	isExported, err := scraper.repo.HasExportedDocument(node.id)
	if err != nil {
		return err
	}

	if !isExported {
		pdf, err := scraper.client.ExportThreadDocument(node.id)
		if err != nil {
			return err
		}
		return scraper.repo.SaveThreadDocument(node.path, node.thread, pdf)
	}
	return nil
}
