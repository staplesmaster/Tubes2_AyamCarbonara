package model

type NodeType int

const (
	ElementNode NodeType = iota
	TextNode
	DocumentNode
)

type DOMNode struct {
	Id         int               `json:"id"`
	Type       NodeType          `json:"-"`
	TagName    string            `json:"tag"`
	Content    string            `json:"content,omitempty"`
	Attributes map[string]string `json:"attrs,omitempty"`
	Children   []*DOMNode        `json:"children"`
	Parent     *DOMNode          `json:"-"`
	Depth      int               `json:"-"`
}

type TraversalRequest struct {
	InputMode string `json:"inputMode"`
	URL       string `json:"url"`
	HTML      string `json:"html"`
	Selector  string `json:"selector"`
	Algorithm string `json:"algorithm"`
	Parallel  bool   `json:"parallel"`
	AllResult bool   `json:"allResult"`
	Limit     int    `json:"limit"`
}

type LCARequest struct {
	InputMode string `json:"inputMode"`
	URL       string `json:"url"`
	HTML      string `json:"html"`
	NodeA     int    `json:"nodeA"`
	NodeB     int    `json:"nodeB"`
}

type LCAStep struct {
	Step          int   `json:"step"`
	NodeA         int   `json:"nodeA"`
	NodeB         int   `json:"nodeB"`
	ActiveNodeIDs []int `json:"activeNodeIds"`
	LCANodeID     *int  `json:"lcaNodeId,omitempty"`
}

type TraversalStep struct {
	Step         int    `json:"step"`
	NodeID       int    `json:"nodeId"`
	Tag          string `json:"tag"`
	Label        string `json:"label"`
	IsMatch      bool   `json:"isMatch"`
	VisitedCount int    `json:"visitedCount"`
	MatchedCount int    `json:"matchedCount"`
}

type TraversalStats struct {
	Visited  int     `json:"visited"`
	Matched  int     `json:"matched"`
	MaxDepth int     `json:"maxDepth"`
	Elapsed  float64 `json:"elapsedMs"`
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

type TraverseResponse struct {
	Success        bool            `json:"success"`
	Tree           *DOMNode        `json:"tree"`
	Steps          []TraversalStep `json:"steps"`
	MatchedNodeIDs []int           `json:"matchedNodeIds"`
	Stats          TraversalStats  `json:"stats"`
	Algorithm      string          `json:"algorithm"`
	Parallel       bool            `json:"parallel"`
	Selector       string          `json:"selector"`
	SourceURL      string          `json:"sourceUrl,omitempty"`
}

type LCAResponse struct {
	Success bool      `json:"success"`
	NodeID  int       `json:"nodeId"`
	Tag     string    `json:"tag"`
	Label   string    `json:"label"`
	Steps   []LCAStep `json:"steps"`
}
