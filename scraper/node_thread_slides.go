package scraper

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

type ThreadSlidesNode struct {
	*ThreadNode
}

func NewThreadSlidesNode(parent *ThreadNode) INode {
	wg, ctx := errgroup.WithContext(parent.ctx)
	return &ThreadSlidesNode{
		ThreadNode: &ThreadNode{
			BaseNode: &BaseNode{
				logger: logrus.WithField("module", NodeTypes.ThreadSlides).
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

func (node *ThreadSlidesNode) ID() string {
	return fmt.Sprintf("thread:%s:slides [%s]", node.id, node.path)
}

func (node *ThreadSlidesNode) Children() []INode {
	return []INode{}
}

func (node *ThreadSlidesNode) Process(scraper *Scraper) error {
	if node.ctx.Err() != nil {
		return nil
	}
	isExported, err := scraper.repo.HasExportedSlides(node.id)
	if err != nil {
		node.logger.Errorln(err)
		return err
	}

	if !isExported {
		pdf, err := scraper.client.ExportThreadSlides(node.id)
		if err != nil {
			node.logger.Errorln(err)
			return err
		}
		return scraper.repo.SaveThreadSlides(node.path, node.thread, pdf)
	} else {
		node.logger.Debugf("already exported")
	}
	return nil
}
