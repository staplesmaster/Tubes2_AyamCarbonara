package selector

import (
	"errors"
	"regexp"
)

type TokenType int

const (
	TAG TokenType = iota
	CLASS
	ID
	ATTR
	COMBINATOR
	UNIVERSAL
)


type Token struct {
	Type  TokenType
	Value string
}

func (t Token) toString() string {
	switch t.Type {
	case CLASS:
		return "." + t.Value
	case ID:
		return "#" + t.Value
	case ATTR:
		return "[" + t.Value + "]"
	default:
		return t.Value
	}
}

func isIdentChar(c byte) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		(c >= '0' && c <= '9') ||
		c == '-' || c == '_'
}

func makeTokens(query string) []Token {
	var result []Token
	i := 0
	prevIsIdent := false
	for i < len(query) {
		if query[i] == ' ' {
			for i < len(query) && query[i] == ' ' {
				i++
			}
			if prevIsIdent && i < len(query) && isIdentChar(query[i]) {
				result = append(result, Token{COMBINATOR, " "})
				prevIsIdent = false
			}
			continue
		}

		switch query[i] {

		case '>', '+', '~':
			result = append(result, Token{COMBINATOR, string(query[i])})
			i++
			prevIsIdent = false

		case '.':
			i++
			start := i
			for i < len(query) && isIdentChar(query[i]) {
				i++
			}
			result = append(result, Token{CLASS, query[start:i]})
			prevIsIdent = true

		case '#':
			i++
			start := i
			for i < len(query) && isIdentChar(query[i]) {
				i++
			}
			result = append(result, Token{ID, query[start:i]})
			prevIsIdent = true

		case '*':
			result = append(result, Token{UNIVERSAL, "*"})
			i++
			prevIsIdent = true

		case '[':
			start := i + 1
			for i < len(query) && query[i] != ']' {
				i++
			}
			result = append(result, Token{ATTR, query[start:i]})
			i++
			prevIsIdent = true

		default:
			start := i
			i++
			for i < len(query) && isIdentChar(query[i]) {
				i++
			}
			result = append(result, Token{TAG, query[start:i]})
			prevIsIdent = true
		}
	}
	return result
}

func getNearTokenString(tokens []Token, i int) string {
	result := "\""
	if (i > 0) {
		result = result + tokens[i-1].toString() + " "
	}
	result = result + tokens[i].toString()
	if (i < len(tokens) - 1) {
		result = result + " " + tokens[i+1].toString() 
	}
	result = result + "\""
	return  result
}

func getTokenError(tokens []Token) error {
	var (
		tagRe   = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9-]*$`)
		classRe = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_-]*$`)
		idRe    = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_-]*$`)
		attrRe  = regexp.MustCompile(`^\s*([a-zA-Z_][a-zA-Z0-9_-]*)(\s*=\s*("[^"]*"|'[^']*'|[^\]\s]+))?\s*$`)
	)
	for i, token := range tokens {
		switch token.Type {
		case TAG:
			if !tagRe.MatchString(token.Value) {
				return errors.New("invalid tagname near " + getNearTokenString(tokens, i))
			}
		case ID:
			if !idRe.MatchString(token.Value) {
				return errors.New("invalid ID near " + getNearTokenString(tokens, i))
			}
		case CLASS:
			if !classRe.MatchString(token.Value) {
				return errors.New("invalid classname near " + getNearTokenString(tokens, i))
			}
		case ATTR:
			if !attrRe.MatchString(token.Value) {
				return errors.New("invalid attribute near " + getNearTokenString(tokens, i))
			}
		}
	}
	return nil
}