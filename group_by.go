package kapacitor

import (
	"sort"

	"github.com/influxdb/kapacitor/models"
	"github.com/influxdb/kapacitor/pipeline"
)

type GroupByNode struct {
	node
	g          *pipeline.GroupByNode
	dimensions []string
}

// Create a new GroupByNode which splits the stream dynamically based on the specified dimensions.
func newGroupByNode(et *ExecutingTask, n *pipeline.GroupByNode) (*GroupByNode, error) {
	gn := &GroupByNode{
		node: node{Node: n, et: et},
		g:    n,
	}
	gn.node.runF = gn.runGroupBy
	gn.dimensions = n.Dimensions
	sort.Strings(gn.dimensions)
	return gn, nil
}

func (g *GroupByNode) runGroupBy() error {
	switch g.Wants() {
	case pipeline.StreamEdge:
		for pt := g.ins[0].NextPoint(); pt != nil; pt = g.ins[0].NextPoint() {
			pt.Group = models.TagsToGroupID(g.dimensions, pt.Tags)
			for _, child := range g.outs {
				err := child.CollectPoint(pt)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
