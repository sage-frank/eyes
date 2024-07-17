package utility

import (
	"fmt"
	"time"

	"github.com/bwmarrin/snowflake"
)

type ISFNode interface {
	GenID() int64
}

type sfNode struct {
	node *snowflake.Node
}

func NewSFNode(startTime string, machineID int64) (ISFNode, error) {
	st, err := time.Parse("2006-01-02", startTime)
	if err != nil {
		return nil, fmt.Errorf("time.Parse err: %w", err)
	}
	snowflake.Epoch = st.UnixNano() / 1000000
	node, err := snowflake.NewNode(machineID)
	if err != nil {
		return nil, fmt.Errorf("snowflake.NewNode(machineID), err: %w", err)
	}

	return &sfNode{
		node: node,
	}, nil
}

func (n *sfNode) GenID() int64 {
	return n.node.Generate().Int64()
}
