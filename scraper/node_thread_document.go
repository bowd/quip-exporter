package scraper

import (
	"github.com/bowd/quip-exporter/interfaces"
	"github.com/bowd/quip-exporter/types"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"path"
)

type ThreadDocumentNode struct {
	*ThreadNode
	exists bool
}

func NewThreadDocumentNode(parent *ThreadNode) interfaces.INode {
	wg, ctx := errgroup.WithContext(parent.ctx)
	return &ThreadDocumentNode{
		ThreadNode: &ThreadNode{
			BaseNode: &BaseNode{
				logger: logrus.WithField("module", types.NodeTypes.ThreadDocument).
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

func (node ThreadDocumentNode) Type() types.NodeType {
	return types.NodeTypes.ThreadDocument
}

func (node ThreadDocumentNode) ID() string {
	return node.id + "/doc"
}

func (node ThreadDocumentNode) Path() string {
	return path.Join("data", "docs", node.id+".docx")
}

func (node *ThreadDocumentNode) Children() []interfaces.INode {
	if !node.exists {
		return []interfaces.INode{}
	}
	return []interfaces.INode{
		NewArchiveNode(
			node.path,
			node.id,
			node.thread.Filename()+".docx",
			node,
		),
	}
}

func (node *ThreadDocumentNode) Process(repo interfaces.IRepository, quip interfaces.IQuipClient, search interfaces.ISearchIndex) error {
	if node.ctx.Err() != nil {
		return nil
	}
	isExported, err := repo.NodeExists(node)
	if err != nil {
		return err
	}
	if !isExported {
		if data, err := quip.ExportThreadDocument(node.id); err != nil {
			return err
		} else {
			if err := repo.SaveNodeRaw(node, data); err != nil {
				return err
			} else {
				node.exists = true
			}
		}
	} else {
		node.exists = true
		node.logger.Debugf("already exported")
	}
	return nil
}
