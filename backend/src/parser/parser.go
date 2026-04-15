package parser

import ( "github.com/luis/Tubes2_AyamCarbonara/backend/src/format_token")

type NodeType int

const(
	Element NodeType = iota
	TextNode
	DocumentNode
)

type DOMNode struct{
	Type NodeType
	TagName string
	Content string
	Attributes map[string]string
	Children []*DOMNode
	Parent *DOMNode
	Depth int
	Index int
}



var nodeCounter int

func ParseHTML(rawHTML string) (*DOMNode, error) {
	index := 0
	formatTokens := format_token.GetFormatToken(rawHTML)
	root := &DOMNode{
		Type:    DocumentNode,
		TagName: "document",
		Depth:   0,
	}

	var parseRecursive func(parent *DOMNode, depth int) ([]*DOMNode, error)
	
	tagWOClosure := map[string]bool{
		"img": true, "br": true, "hr": true, "input": true, "link": true, "meta": true,
	}

	parseRecursive = func(parent *DOMNode, depth int) ([]*DOMNode, error) {
		var children []*DOMNode

		for index < len(formatTokens) {
			token := formatTokens[index]

			if token.Kind == format_token.FormatClosingTag {
				if token.TagName == parent.TagName {
					index++
					return children, nil
				}
				return children, nil 
			}

			if token.Kind == format_token.FormatText {
				node := &DOMNode{
					Type:    TextNode,
					Content: token.Content,
					Parent:  parent,
					Depth:   depth,
				}
				children = append(children, node)
				index++
			} else if token.Kind == format_token.FormatOpeningTag {
				newNode := &DOMNode{
					Type:       Element,
					TagName:    token.TagName,
					Attributes: token.Attributes,
					Parent:     parent,
					Depth:      depth,
				}
				index++

				if !tagWOClosure[newNode.TagName] && token.Content != "self-closing" {
					childList, err := parseRecursive(newNode, depth+1)
					if err != nil {
						return nil, err
					}
					newNode.Children = childList
				}
				children = append(children, newNode)
			}
		}
		return children, nil
	}

	children, err := parseRecursive(root, 1)
	root.Children = children
	return root, err
}

// Helper: hitung max depth tree
func MaxDepth(root *DOMNode) int {
    if root == nil || len(root.Children) == 0 {
        return root.Depth
    }
    max := 0
    for _, child := range root.Children {
        d := MaxDepth(child)
        if d > max {
            max = d
        }
    }
    return max
}

// Helper: hitung total node
func CountNodes(root *DOMNode) int {
    if root == nil {
        return 0
    }
    count := 1
    for _, child := range root.Children {
        count += CountNodes(child)
    }
    return count
}