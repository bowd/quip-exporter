package scraper

import (
	"fmt"
	"golang.org/x/sync/errgroup"
)

type ThreadHTMLNode struct {
	*ThreadNode
}

func NewThreadHTMLNode(parent *ThreadNode) INode {
	wg, ctx := errgroup.WithContext(parent.ctx)
	return &ThreadHTMLNode{
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

func (node *ThreadHTMLNode) ID() string {
	return fmt.Sprintf("thread:%s:html [%s]", node.id, node.path)
}

func (node *ThreadHTMLNode) Children() []INode {
	return []INode{}
}

func (node *ThreadHTMLNode) Process(scraper *Scraper) error {
	scraper.logger.Debugf("Handling thread:%s:HTML [%s/%s]", node.id, node.path, node.thread.Filename())
	if node.ctx.Err() != nil {
		return nil
	}
	isExported, err := scraper.repo.HasExportedHTML(node.id)
	if err != nil {
		return err
	}

	if !isExported {
		return scraper.repo.SaveThreadHTML(node.path, node.thread)
	}
	return nil
}
