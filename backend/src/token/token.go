package token

import (
	 "unicode"
	 "fmt"
)

type TokenKind int
const(
	TokenOpen TokenKind = iota
	TokenOpenSlash
	TokenClose
	TokenSelfClose
	TokenEquals
	TokenWord
	TokenText
)

type Token struct{
	Kind TokenKind
	Value string
}

type State int
const(
	StateOutside State = iota
	StateSeenOpenBracket
	StateInsideTag
	StateSeenSlash
	StateInsideQuote
	StateError
)

func Tokenize(rawHTML string) ([]Token, error){
	var res []Token
	var currentWord []rune
	state := StateOutside
	var quoteChar rune

	runes := []rune(rawHTML)
	for i := 0; i < len(runes); i++{
		cc := runes[i]

		switch state{
		case StateOutside:
			if cc == '<'{
				if len(currentWord) > 0{
					res = append(res, Token{TokenText, string(currentWord)})
					currentWord = nil
				}
				state = StateSeenOpenBracket
			}else{
					currentWord = append(currentWord, cc)
			}	
		case StateSeenOpenBracket: 
			if cc == '/'{
				res = append(res, Token{TokenOpenSlash, "</"})
				state = StateInsideTag
			}else{
				res = append(res, Token{TokenOpen, "<"})
				state = StateInsideTag
				i--
			}
		case StateInsideTag:
			if unicode.IsSpace(cc) {
                if len(currentWord) > 0 {
                    res = append(res, Token{TokenWord, string(currentWord)})
                    currentWord = nil
                }
            } else if cc == '>' {
                if len(currentWord) > 0 {
                    res = append(res, Token{TokenWord, string(currentWord)})
                    currentWord = nil
                }
                res = append(res, Token{TokenClose, ">"})
                state = StateOutside
            } else if cc == '/' {
                state = StateSeenSlash
            } else if cc == '=' {
                if len(currentWord) > 0 {
                    res = append(res, Token{TokenWord, string(currentWord)})
                    currentWord = nil
                }
                res = append(res, Token{TokenEquals, "="})
            } else if cc == '"' || cc == '\'' {
                quoteChar = cc
                state = StateInsideQuote
            } else if cc == '<'{
				state = StateError
			}else {
                currentWord = append(currentWord, cc)
            }	
		case StateInsideQuote:
			if cc == quoteChar{
				res = append(res, Token{TokenWord, string(currentWord)})
				currentWord = nil
				state = StateInsideTag
			}else{
				currentWord = append(currentWord, cc)
			}
		case StateSeenSlash:
			if cc == '>'{
				res = append(res, Token{TokenSelfClose, "/>"})
				state = StateOutside
			}else{
				state = StateInsideTag
				i--
			}
		case StateError:
			return res, fmt.Errorf("File HTML tidak sesuai format yang telah ditentukan. Kirim file yang valid!")
		}
	}

	if state == StateInsideQuote || state == StateSeenSlash || state == StateInsideTag || state == StateSeenOpenBracket{
		return res, fmt.Errorf("File HTML tidak sesuai format yang telah ditentukan. Kirim file yang valid!")
	}

	return res, nil
}
