package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/bowd/quip-exporter/interfaces"
	"github.com/bowd/quip-exporter/types"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"path"
	"strings"
)

type ThreadIndexNode struct {
	*ThreadNode
	exists bool
}

func NewThreadIndexNode(parent *ThreadNode) interfaces.INode {
	wg, ctx := errgroup.WithContext(parent.ctx)
	return &ThreadIndexNode{
		ThreadNode: &ThreadNode{
			BaseNode: &BaseNode{
				logger: logrus.WithField("module", types.NodeTypes.ThreadIndex).
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

func (node ThreadIndexNode) Type() types.NodeType {
	return types.NodeTypes.ThreadIndex
}

func (node *ThreadIndexNode) ID() string {
	return node.id + "/index"
}

func (node ThreadIndexNode) Path() string {
	return path.Join("data", "index", node.id)
}

func (node *ThreadIndexNode) Children() []interfaces.INode {
	return []interfaces.INode{}
}

func (node *ThreadIndexNode) Process(repo interfaces.IRepository, quip interfaces.IQuipClient, search interfaces.ISearchIndex) error {
	if node.ctx.Err() != nil {
		return nil
	}

	if indexed, err := search.IsIndexed(node.thread.Thread.ID); err != nil {
		return err
	} else if indexed {
		node.logger.Debug("Skipping already indexed")
		return nil
	}

	if doc, err := goquery.NewDocumentFromReader(strings.NewReader(node.thread.HTML)); err != nil {
		return err
	} else {
		return search.IndexThread(
			node.thread.Thread.ID,
			node.thread.Thread.Title,
			doc.Text(),
		)
	}

}
