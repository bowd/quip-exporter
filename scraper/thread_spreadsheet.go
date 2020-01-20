package scraper

import (
	"fmt"
	"golang.org/x/sync/errgroup"
)

type ThreadSpreadsheetNode struct {
	*ThreadNode
}

func NewThreadSpreadsheetNode(parent *ThreadNode) INode {
	wg, ctx := errgroup.WithContext(parent.ctx)
	return &ThreadSpreadsheetNode{
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

func (node *ThreadSpreadsheetNode) ID() string {
	return fmt.Sprintf("thread:%s:spreadsheet [%s]", node.id, node.path)
}

func (node *ThreadSpreadsheetNode) Children() []INode {
	return []INode{}
}

func (node *ThreadSpreadsheetNode) Process(scraper *Scraper) error {
	scraper.logger.Debugf("Handling thread:%s:XLS [%s/%s]", node.id, node.path, node.thread.Filename())
	if node.ctx.Err() != nil {
		return nil
	}
	isExported, err := scraper.repo.HasExportedSpreadsheet(node.id)
	if err != nil {
		return err
	}

	if !isExported {
		pdf, err := scraper.client.ExportThreadSpreadsheet(node.id)
		if err != nil {
			return err
		}
		return scraper.repo.SaveThreadSpreadsheet(node.path, node.thread, pdf)
	}
	return nil
}
