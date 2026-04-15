package selector

import (
	"strings"

	"github.com/luis/Tubes2_AyamCarbonara/backend/src/model"
)

type Selector func(node *model.DOMNode) bool

func TagSelector(tag string) Selector {
	return func(node *model.DOMNode) bool {
		return node.TagName == tag
	}
}

func ClassSelector(class string) Selector {
	return func(node *model.DOMNode) bool {
		if classGroup, ok := node.Attributes["class"]; ok {
			classes := strings.Split(classGroup, " ")
			for _, c := range classes {
				if c == class {
					return true
				}
			}
		}
		return false
	}
}

func IDSelector(id string) Selector {
	return func(node *model.DOMNode) bool {
		if nodeId, ok := node.Attributes["id"]; ok {
			return nodeId == id
		}
		return false
	}
}

func UniversalSelector() Selector {
	return func(node *model.DOMNode) bool {
		return true
	}
}

func And(selectors ...Selector) Selector {
	return func(node *model.DOMNode) bool {
		for _, m := range selectors {
			if !m(node) {
				return false
			}
		}
		return true
	}
}

func AttributeSelector(key, value string) Selector {
	return func(node *model.DOMNode) bool {
		if attrVal, ok := node.Attributes[key]; ok {
			return attrVal == value
		} 
		return false
	}
}

func Descendant(parentSelector, childSelector Selector) Selector {
	return func(node *model.DOMNode) bool {
		if !childSelector(node) {
			return false
		}
		for p := node.Parent; p != nil; p = p.Parent {
			if parentSelector(p) {
				return true
			}
		}
		return false
	}
}

func Child(parentSelector, childSelector Selector) Selector {
	return func(node *model.DOMNode) bool {
		return node.Parent != nil &&
			parentSelector(node.Parent) &&
			childSelector(node)
	}
}

func AdjacentSibling(prevSelector, currSelector Selector) Selector {
	return func(node *model.DOMNode) bool {
		if !currSelector(node) || node.Parent == nil {
			return false
		}
		siblings := node.Parent.Children
		for i := 1; i < len(siblings); i++ {
			if siblings[i] == node {
				return prevSelector(siblings[i-1])
			}
		}
		return false
	}
}

func GeneralSibling(prevSelector, currSelector Selector) Selector {
	return func(node *model.DOMNode) bool {
		if !currSelector(node) || node.Parent == nil {
			return false
		}
		siblings := node.Parent.Children
		for i := 0; i < len(siblings); i++ {
			if siblings[i] == node {
				for j := 0; j < i; j++ {
					if prevSelector(siblings[j]) {
						return true
					}
				}
			}
		}
		return false
	}
}