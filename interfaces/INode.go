package interfaces

import "github.com/bowd/quip-exporter/types"

type INode interface {
	Go(func() error)
	Wait() error
	Children() []INode
	Process(IRepository, IQuipClient, ISearchIndex) error
	Type() types.NodeType
	ID() string
	Path() string
}
