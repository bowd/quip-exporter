package scraper

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"regexp"
)

type ThreadHTMLNode struct {
	*ThreadNode
}

var blobRegexp *regexp.Regexp = regexp.MustCompile("blob/([0-9a-zA-Z_-]+)/([0-9a-zA-z_-]+)")

func NewThreadHTMLNode(parent *ThreadNode) INode {
	wg, ctx := errgroup.WithContext(parent.ctx)
	return &ThreadHTMLNode{
		ThreadNode: &ThreadNode{
			BaseNode: &BaseNode{
				logger: logrus.WithField("module", NodeTypes.ThreadHTML).
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

func (node *ThreadHTMLNode) ID() string {
	return fmt.Sprintf("thread:%s:html [%s]", node.id, node.path)
}

func (node *ThreadHTMLNode) Children() []INode {
	return []INode{}
	// matches := blobRegexp.FindAllStringSubmatch(node.thread.HTML, -1)
	// children := make([]INode, 0, 0)
	// for _, match := range matches {
	// 	if match[1] == node.id {
	// 		children = append(children, NewBlobNode(node.ThreadNode, match[2]))
	// 	}
	// }
	// return children
}

func (node *ThreadHTMLNode) Process(scraper *Scraper) error {
	if node.ctx.Err() != nil {
		return nil
	}
	isExported, err := scraper.repo.HasExportedHTML(node.id)
	if err != nil {
		node.logger.Errorln(err)
		return err
	}

	if !isExported {
		err := scraper.repo.SaveThreadHTML(node.path, node.thread)
		if err != nil {
			node.logger.Errorln(err)
		}
		return err
	} else {
		node.logger.Debugf("already exported")
	}

	return nil
}
