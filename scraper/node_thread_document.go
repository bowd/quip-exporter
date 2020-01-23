package scraper

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

type ThreadDocumentNode struct {
	*ThreadNode
}

func NewThreadDocumentNode(parent *ThreadNode) INode {
	wg, ctx := errgroup.WithContext(parent.ctx)
	return &ThreadDocumentNode{
		ThreadNode: &ThreadNode{
			BaseNode: &BaseNode{
				logger: logrus.WithField("module", NodeTypes.ThreadDocument).
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
		node.logger.Errorln(err)
		return err
	}

	if !isExported {
		pdf, err := scraper.client.ExportThreadDocument(node.id)
		if err != nil {
			node.logger.Errorln(err)
			return err
		}
		return scraper.repo.SaveThreadDocument(node.path, node.thread, pdf)
	} else {
		node.logger.Debugf("already exported")
	}
	return nil
}
