package scraper

import (
	"github.com/bowd/quip-exporter/interfaces"
	"github.com/bowd/quip-exporter/types"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"path"
	"regexp"
)

type ThreadHTMLNode struct {
	*ThreadNode
	exists bool
}

var blobRegexp *regexp.Regexp = regexp.MustCompile("blob/([0-9a-zA-Z_-]+)/([0-9a-zA-z_-]+)")

func NewThreadHTMLNode(parent *ThreadNode) interfaces.INode {
	wg, ctx := errgroup.WithContext(parent.ctx)
	return &ThreadHTMLNode{
		ThreadNode: &ThreadNode{
			BaseNode: &BaseNode{
				logger: logrus.WithField("module", types.NodeTypes.ThreadHTML).
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

func (node ThreadHTMLNode) Type() types.NodeType {
	return types.NodeTypes.ThreadHTML
}

func (node *ThreadHTMLNode) ID() string {
	return node.id + "/html"
}

func (node ThreadHTMLNode) Path() string {
	return path.Join("data", "html", node.id+".html")
}

func (node *ThreadHTMLNode) Children() []interfaces.INode {
	if !node.exists {
		return []interfaces.INode{}
	}
	matches := blobRegexp.FindAllStringSubmatch(node.thread.HTML, -1)
	children := make([]interfaces.INode, 0, 0)
	for _, match := range matches {
		if match[1] == node.id {
			children = append(children, NewBlobNode(node.ThreadNode, match[2]))
		}
	}
	children = append(
		children,
		NewArchiveNode(
			node.path,
			node.id,
			node.thread.Filename()+".html",
			node,
		),
	)
	return children
}

func (node *ThreadHTMLNode) Process(repo interfaces.IRepository, quip interfaces.IQuipClient) error {
	if node.ctx.Err() != nil {
		return nil
	}
	isExported, err := repo.NodeExists(node)
	if err != nil {
		node.logger.Errorln(err)
		return err
	}

	if !isExported {
		if err := repo.SaveNodeRaw(node, []byte(node.thread.HTML)); err != nil {
			node.logger.Errorln(err)
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
