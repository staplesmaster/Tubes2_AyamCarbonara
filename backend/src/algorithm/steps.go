package algorithm

import (
	"fmt"

	"github.com/luis/Tubes2_AyamCarbonara/backend/src/model"
	"github.com/luis/Tubes2_AyamCarbonara/backend/src/selector"
)

func visitStep(
	n *model.DOMNode,
	sel selector.Selector,
	step *int,
	visited *int,
	matched *int,
	maxDepth *int,
	steps *[]model.TraversalStep,
	matchedIDs *[]int,
) {
	(*visited)++
	(*step)++

	if n.Depth > *maxDepth {
		*maxDepth = n.Depth
	}

	isMatch := sel(n)
	if isMatch {
		(*matched)++
		*matchedIDs = append(*matchedIDs, n.Id)
	}

	*steps = append(*steps, model.TraversalStep{
		Step:         *step,
		NodeID:       n.Id,
		Tag:          n.TagName,
		Label:        fmt.Sprintf("%s#%d", n.TagName, n.Id),
		IsMatch:      isMatch,
		VisitedCount: *visited,
		MatchedCount: *matched,
	})
}
