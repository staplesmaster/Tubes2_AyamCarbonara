package selector

import (
	"errors"
	"strings"
)

func buildSelector(tokens []Token) (Selector, error) {
	var compounds []Selector
	var combinators []Token
	allowCombinator := false
	allowTagUniversal := true
	currentCompound := UniversalSelector()
	for _, token := range tokens {
		switch token.Type {
		case COMBINATOR:
			if !allowCombinator {
				return nil, errors.New("expected identifier but found combinator")
			}
			compounds = append(compounds, currentCompound)
			currentCompound = UniversalSelector()
			allowTagUniversal = true

			combinators = append(combinators, token)
			allowCombinator = false
		case UNIVERSAL:
			if !allowTagUniversal {
				return nil, errors.New("unexpected occurrence of universal selector")
			}
			allowTagUniversal = false
			allowCombinator = true
		case TAG:
			if !allowTagUniversal {
				return nil, errors.New("unexpected occurrence of tag selector")
			}
			currentCompound = TagSelector(token.Value)
			allowTagUniversal = false
			allowCombinator = true
		case CLASS:
			currentCompound = And(ClassSelector(token.Value), currentCompound)
			allowTagUniversal = false
			allowCombinator = true
		case ID:
			currentCompound = And(IDSelector(token.Value), currentCompound)
			allowTagUniversal = false
			allowCombinator = true
		case ATTR:
			vals := strings.Split(token.Value, "=")
			if len(vals) > 1 {
				currentCompound = And(MatchAttributeSelector(vals[0], vals[1]), currentCompound)
			} else {
				currentCompound = And(HasAttributeSelector(vals[0]), currentCompound)
			}
			allowTagUniversal = false
			allowCombinator = true
		}
	}
	if !allowCombinator {
		return nil, errors.New("hanging last combinator")
	}
	compounds = append(compounds, currentCompound)

	if len(combinators) == 0 {
		return compounds[0], nil
	}

	if len(compounds) != len(combinators)+1 {
		return nil, errors.New("invalid selector structure")
	}

	finalSelector := compounds[0]

	for j := 0; j < len(combinators); j++ {
		switch combinators[j].Value {
		case " ":
			finalSelector = Descendant(finalSelector, compounds[j+1])
		case ">":
			finalSelector = Child(finalSelector, compounds[j+1])
		case "+":
			finalSelector = AdjacentSibling(finalSelector, compounds[j+1])
		case "~":
			finalSelector = GeneralSibling(finalSelector, compounds[j+1])
		}
	}

	return finalSelector, nil
}

func StringToSelector(query string) (Selector, error) {
	tokens := makeTokens(query)
	if err := getTokenError(tokens); err != nil {
		return nil, err
	}
	return buildSelector(tokens)
}