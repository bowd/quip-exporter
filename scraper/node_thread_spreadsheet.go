package scraper

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

type ThreadSpreadsheetNode struct {
	*ThreadNode
}

func NewThreadSpreadsheetNode(parent *ThreadNode) INode {
	wg, ctx := errgroup.WithContext(parent.ctx)
	return &ThreadSpreadsheetNode{
		ThreadNode: &ThreadNode{
			BaseNode: &BaseNode{
				logger: logrus.WithField("module", NodeTypes.ThreadSpreadsheet).
					WithField("id", parent.id).
					WithField("path", parent.path),
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
	if node.ctx.Err() != nil {
		return nil
	}
	isExported, err := scraper.repo.HasExportedSpreadsheet(node.id)
	if err != nil {
		node.logger.Errorln(err)
		return err
	}

	if !isExported {
		pdf, err := scraper.client.ExportThreadSpreadsheet(node.id)
		if err != nil {
			node.logger.Errorln(err)
			return err
		}
		return scraper.repo.SaveThreadSpreadsheet(node.path, node.thread, pdf)
	} else {
		node.logger.Debugf("already exported")
	}
	return nil
}
