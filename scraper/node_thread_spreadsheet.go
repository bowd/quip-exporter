package scraper

import (
	"github.com/bowd/quip-exporter/interfaces"
	"github.com/bowd/quip-exporter/types"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"path"
)

type ThreadSpreadsheetNode struct {
	*ThreadNode
	exists bool
}

func NewThreadSpreadsheetNode(parent *ThreadNode) interfaces.INode {
	wg, ctx := errgroup.WithContext(parent.ctx)
	return &ThreadSpreadsheetNode{
		ThreadNode: &ThreadNode{
			BaseNode: &BaseNode{
				logger: logrus.WithField("module", types.NodeTypes.ThreadSpreadsheet).
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

func (node ThreadSpreadsheetNode) Type() types.NodeType {
	return types.NodeTypes.ThreadSpreadsheet
}

func (node *ThreadSpreadsheetNode) ID() string {
	return node.id + "/xls"
}

func (node ThreadSpreadsheetNode) Path() string {
	return path.Join("data", "xls", node.id+".xlsx")
}

func (node *ThreadSpreadsheetNode) Children() []interfaces.INode {
	if !node.exists {
		return []interfaces.INode{}
	}
	return []interfaces.INode{
		NewArchiveNode(
			node.path,
			node.id,
			node.thread.Filename()+".xlsx",
			node,
		),
	}
}

func (node *ThreadSpreadsheetNode) Process(repo interfaces.IRepository, quip interfaces.IQuipClient, search interfaces.ISearchIndex) error {
	if node.ctx.Err() != nil {
		return nil
	}
	isExported, err := repo.NodeExists(node)
	if err != nil {
		return err
	}

	if !isExported {
		data, err := quip.ExportThreadSpreadsheet(node.id)
		if err != nil {
			return err
		}
		if err := repo.SaveNodeRaw(node, data); err != nil {
			return err
		} else {
			node.exists = true
		}
	} else {
		node.exists = true
		node.logger.Debugf("already exported")
	}
	return nil
}
