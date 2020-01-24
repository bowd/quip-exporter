package interfaces

import "github.com/bowd/quip-exporter/types"

type INode interface {
	Go(func() error)
	Wait() error
	Children() []INode
	Process(repo IRepository, quip IQuipClient) error
	Type() types.NodeType
	ID() string
	Path() string
}
