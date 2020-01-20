package scraper

import (
	"fmt"
	"golang.org/x/sync/errgroup"
)

type ThreadSlidesNode struct {
	*ThreadNode
}

func NewThreadSlidesNode(parent *ThreadNode) INode {
	wg, ctx := errgroup.WithContext(parent.ctx)
	return &ThreadSlidesNode{
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

func (node *ThreadSlidesNode) ID() string {
	return fmt.Sprintf("thread:%s:slides [%s]", node.id, node.path)
}

func (node *ThreadSlidesNode) Children() []INode {
	return []INode{}
}

func (node *ThreadSlidesNode) Process(scraper *Scraper) error {
	scraper.logger.Debugf("Handling thread:%s:PDF [%s/%s]", node.id, node.path, node.thread.Filename())
	if node.ctx.Err() != nil {
		return nil
	}
	isExported, err := scraper.repo.HasExportedSlides(node.id)
	if err != nil {
		return err
	}

	if !isExported {
		pdf, err := scraper.client.ExportThreadSlides(node.id)
		if err != nil {
			return err
		}
		return scraper.repo.SaveThreadSlides(node.path, node.thread, pdf)
	}
	return nil
}
