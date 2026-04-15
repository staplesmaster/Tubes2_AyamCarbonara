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