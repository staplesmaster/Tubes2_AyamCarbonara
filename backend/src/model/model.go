package model

type DOMNode struct {
	TagName    string
	Attributes map[string]string
	Children   []*DOMNode
	Parent     *DOMNode
	Depth      int
}
