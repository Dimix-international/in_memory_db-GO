package service

import (
	"strings"

	"github.com/Dimix-international/in_memory_db-GO/internal/models"
)

const (
	startState = iota
	letterOrPunctuationState
	whiteSpaceState
)

type ParserService struct {
	state  int
	tokens []string
	sb     strings.Builder
}

func NewParserService() *ParserService {
	return &ParserService{state: startState}
}

func (p *ParserService) Parse(query string) ([]string, error) {
	for i := range query {
		switch p.state {
		case startState:
			if !isLetterOrPunctuationSymbol(query[i]) {
				return nil, models.ErrParsing
			}
			p.sb.WriteByte(query[i])
			p.state = letterOrPunctuationState
		case letterOrPunctuationState:
			if isSpaceSymbol(query[i]) {
				p.tokens = append(p.tokens, p.sb.String())
				p.sb.Reset()
				p.state = whiteSpaceState
				break
			}
			if !isLetterOrPunctuationSymbol(query[i]) {
				return nil, models.ErrParsing
			}
			p.sb.WriteByte(query[i])
		case whiteSpaceState:
			if isSpaceSymbol(query[i]) {
				continue
			}
			if !isLetterOrPunctuationSymbol(query[i]) {
				return nil, models.ErrParsing
			}

			p.sb.WriteByte(query[i])
			p.state = letterOrPunctuationState
		}
	}

	if p.state == letterOrPunctuationState {
		p.tokens = append(p.tokens, p.sb.String())
	}

	return p.tokens, nil
}

func isSpaceSymbol(symbol byte) bool {
	return symbol == '\t' || symbol == '\n' || symbol == ' '
}

func isLetterOrPunctuationSymbol(symbol byte) bool {
	return (symbol >= 'a' && symbol <= 'z') ||
		(symbol >= 'A' && symbol <= 'Z') ||
		(symbol >= '0' && symbol <= '9') ||
		(symbol == '*') ||
		(symbol == '/') ||
		(symbol == '_')
}
