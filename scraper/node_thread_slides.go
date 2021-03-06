package scraper

import (
	"github.com/bowd/quip-exporter/interfaces"
	"github.com/bowd/quip-exporter/types"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"path"
)

type ThreadSlidesNode struct {
	*ThreadNode
	exists bool
}

func NewThreadSlidesNode(parent *ThreadNode) interfaces.INode {
	wg, ctx := errgroup.WithContext(parent.ctx)
	return &ThreadSlidesNode{
		ThreadNode: &ThreadNode{
			BaseNode: &BaseNode{
				logger: logrus.WithField("module", types.NodeTypes.ThreadSlides).
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

func (node ThreadSlidesNode) Type() types.NodeType {
	return types.NodeTypes.ThreadSlides
}

func (node *ThreadSlidesNode) ID() string {
	return node.id + "/pdf"
}

func (node ThreadSlidesNode) Path() string {
	return path.Join("data", "pdf", node.id+".pdf")
}

func (node *ThreadSlidesNode) Children() []interfaces.INode {
	if !node.exists {
		return []interfaces.INode{}
	}
	return []interfaces.INode{
		NewArchiveNode(
			node.path,
			node.id,
			node.thread.Filename()+".pdf",
			node,
		),
	}
}

func (node *ThreadSlidesNode) Process(repo interfaces.IRepository, quip interfaces.IQuipClient, search interfaces.ISearchIndex) error {
	if node.ctx.Err() != nil {
		return nil
	}
	isExported, err := repo.NodeExists(node)
	if err != nil {
		return err
	}

	if !isExported {
		data, err := quip.ExportThreadSlides(node.id)
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
