package model

type NodeType int
const (
    ElementNode NodeType = iota
    TextNode
    DocumentNode
)

type DOMNode struct {
    Id         int
    Type       NodeType
    TagName    string
    Content    string 
    Attributes map[string]string
    Children   []*DOMNode
    Parent     *DOMNode
    Depth      int
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

type APINode struct {
	ID       int               `json:"id"`
	Tag      string            `json:"tag"`
	Attrs    map[string]string `json:"attrs"`
	Children []APINode         `json:"children"`
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
	Visited  int   `json:"visited"`
	Matched  int   `json:"matched"`
	MaxDepth int   `json:"maxDepth"`
	Elapsed  int64 `json:"elapsedMs"`
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

type TraverseResponse struct {
	Success        bool            `json:"success"`
	Tree           APINode         `json:"tree"`
	Steps          []TraversalStep `json:"steps"`
	MatchedNodeIDs []int           `json:"matchedNodeIds"`
	Stats          TraversalStats  `json:"stats"`
	Algorithm      string          `json:"algorithm"`
	Parallel       bool            `json:"parallel"`
	Selector       string          `json:"selector"`
	SourceURL      string          `json:"sourceUrl,omitempty"`
}

type LCAResponse struct {
	Success bool   `json:"success"`
	NodeID  int    `json:"nodeId"`
	Tag     string `json:"tag"`
	Label   string `json:"label"`
}