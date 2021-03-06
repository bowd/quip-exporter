package scraper

import (
	"context"
	"github.com/bowd/quip-exporter/interfaces"
	"github.com/bowd/quip-exporter/types"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"path"
)

type ArchiveNode struct {
	*BaseNode
	source   interfaces.INode
	filename string
}

func NewArchiveNode(path, id, filename string, source interfaces.INode) interfaces.INode {
	wg, _ := errgroup.WithContext(context.Background())
	return &ArchiveNode{
		BaseNode: &BaseNode{
			id:   id,
			path: path,
			wg:   wg,
			logger: logrus.WithField("module", types.NodeTypes.Archive).
				WithField("source", source.Type).
				WithField("filename", filename).
				WithField("path", path),
		},
		source:   source,
		filename: filename,
	}
}

func (node ArchiveNode) Path() types.NodeType {
	return path.Join("archive", node.path, node.filename)
}

func (node ArchiveNode) Type() types.NodeType {
	return types.NodeTypes.Archive
}

func (node ArchiveNode) Children() []interfaces.INode {
	return []interfaces.INode{}
}

func (node ArchiveNode) Process(repo interfaces.IRepository, quip interfaces.IQuipClient, search interfaces.ISearchIndex) error {
	exists, err := repo.NodeExists(node)
	if err != nil {
		return err
	}

	if !exists {
		err := repo.MakeArchiveCopy(node.source, node)
		if err != nil {
			return err
		}
	} else {
		node.logger.Debugf("found in repo")
	}
	return nil
}
