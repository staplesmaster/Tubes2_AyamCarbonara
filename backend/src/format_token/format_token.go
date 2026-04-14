package format_token


import("github.com/luis/Tubes2_AyamCarbonara/backend/src/token"
"fmt")

type FormatTokenKind int

const(
	FormatOpeningTag FormatTokenKind = iota
	FormatClosingTag
	FormatText
)

type FormatToken struct{
	Kind FormatTokenKind
	TagName string
	Attributes map[string]string
	Content string
}

func GetFormatToken(rawHTML string) []FormatToken {
	tokens, err := token.Tokenize(rawHTML)
	if err != nil{
		fmt.Printf("Terjadi error saat tokenize: %w", err)
	}
	var formatTokens []FormatToken
	i := 0

	for i < len(tokens) {
		t := tokens[i]

		switch t.Kind {
		case token.TokenText:
			formatTokens = append(formatTokens, FormatToken{Kind: FormatText, Content: t.Value})
			i++

		case token.TokenOpen:
			i++ 
			if i < len(tokens) && tokens[i].Kind == token.TokenWord {
				tag := FormatToken{
					Kind:       FormatOpeningTag,
					TagName:    tokens[i].Value,
					Attributes: make(map[string]string),
				}
				i++
				for i < len(tokens) && tokens[i].Kind != token.TokenClose && tokens[i].Kind != token.TokenSelfClose {
					if tokens[i].Kind == token.TokenWord {
						key := tokens[i].Value
						i++
						if i < len(tokens) && tokens[i].Kind == token.TokenEquals {
							i++
							if i < len(tokens) && tokens[i].Kind == token.TokenWord {
								tag.Attributes[key] = tokens[i].Value
								i++
							}
						} else {
							tag.Attributes[key] = "true" 
						}
					} else {
						i++
					}
				}
				
				if i < len(tokens) && tokens[i].Kind == token.TokenSelfClose {
					tag.Content = "self-closing" 
				}
				formatTokens = append(formatTokens, tag)
				i++ 
			}

		case token.TokenOpenSlash:
			i++ 
			if i < len(tokens) && tokens[i].Kind == token.TokenWord {
				formatTokens = append(formatTokens, FormatToken{
					Kind:    FormatClosingTag,
					TagName: tokens[i].Value,
				})
				i++
			}
			for i < len(tokens) && tokens[i].Kind != token.TokenClose {
				i++
			}
			i++
		default:
			i++
		}
	}
	return formatTokens
}