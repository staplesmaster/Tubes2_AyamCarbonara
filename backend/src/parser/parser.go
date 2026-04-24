package parser

import ( "github.com/luis/Tubes2_AyamCarbonara/backend/src/format_token"
"github.com/luis/Tubes2_AyamCarbonara/backend/src/model")



func ParseHTML(rawHTML string) (*model.DOMNode, error) {
	index := 0
	nodeCounter := 0
	formatTokens := format_token.GetFormatToken(rawHTML)
	root := &model.DOMNode{
		Id:    nodeCounter,
		TagName: "document",
		Depth:   0,
	}
	nodeCounter++

	var parseRecursive func(parent *model.DOMNode, depth int) ([]*model.DOMNode, error)
	
	tagWOClosure := map[string]bool{
		"img": true, "br": true, "hr": true, "input": true, "link": true, "meta": true,
	}

	parseRecursive = func(parent *model.DOMNode, depth int) ([]*model.DOMNode, error) {
		var children []*model.DOMNode

		for index < len(formatTokens) {
			token := formatTokens[index]

			if token.Kind == format_token.FormatClosingTag {
				if token.TagName == parent.TagName {
					index++
					return children, nil
				}
				index++
				continue // buat kasus <div> ayam </..> </div>
			}

			if token.Kind == format_token.FormatText {
				node := &model.DOMNode{
					Id:    nodeCounter,
					Type: model.TextNode,
					Content: token.Content,
					Parent:  parent,
					Depth:   depth,
				}
				nodeCounter++
				children = append(children, node)
				index++
			} else if token.Kind == format_token.FormatOpeningTag {
				if token.TagName == "!--" {
					index++
					continue
				}
				newNode := &model.DOMNode{
					Id:       nodeCounter,
					Type: model.ElementNode,
					TagName:    token.TagName,
					Attributes: token.Attributes,
					Parent:     parent,
					Depth:      depth,
				}
				nodeCounter++
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

