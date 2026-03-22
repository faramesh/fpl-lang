package parser

import (
	"fmt"

	"github.com/faramesh/fpl-lang/reference/go/fpl/lexer"
	"github.com/faramesh/fpl-lang/reference/go/fpl/token"
)

type Document struct {
	Agents []Agent `json:"agents"`
}

type Agent struct {
	Name    string `json:"name"`
	Default string `json:"default,omitempty"`
	Rules   []Rule `json:"rules,omitempty"`
}

type Rule struct {
	Effect    string `json:"effect"`
	Tool      string `json:"tool"`
	Condition *Expr  `json:"condition,omitempty"`
	Notify    string `json:"notify,omitempty"`
	Reason    string `json:"reason,omitempty"`
	Reeval    *bool  `json:"reeval,omitempty"`
}

type Expr struct {
	Kind  string `json:"kind"`
	Op    string `json:"op,omitempty"`
	Value string `json:"value,omitempty"`
	Left  *Expr  `json:"left,omitempty"`
	Right *Expr  `json:"right,omitempty"`
}

type Parser struct {
	l         *lexer.Lexer
	curToken  token.Token
	peekToken token.Token
}

func New(input string) *Parser {
	l := lexer.New(input)
	p := &Parser{l: l}
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) ParseDocument() (*Document, error) {
	doc := &Document{}

	for p.curToken.Type != token.EOF {
		p.consumeNewlines()
		if p.curToken.Type == token.EOF {
			break
		}

		if p.curToken.Type != token.Agent {
			return nil, fmt.Errorf("unexpected token %q at %d:%d", p.curToken.Literal, p.curToken.Line, p.curToken.Column)
		}

		ag, err := p.parseAgent()
		if err != nil {
			return nil, err
		}
		doc.Agents = append(doc.Agents, *ag)
	}

	if len(doc.Agents) == 0 {
		return nil, fmt.Errorf("document must contain at least one agent block")
	}

	return doc, nil
}

func (p *Parser) parseAgent() (*Agent, error) {
	if !p.expectPeek(token.Ident) {
		return nil, p.unexpected("expected agent name after agent")
	}
	agent := &Agent{Name: p.curToken.Literal}

	if !p.expectPeek(token.LBrace) {
		return nil, p.unexpected("expected '{' after agent name")
	}
	p.nextToken()

	for p.curToken.Type != token.RBrace {
		if p.curToken.Type == token.EOF {
			return nil, fmt.Errorf("unterminated agent block %q", agent.Name)
		}
		if p.curToken.Type == token.Newline {
			p.nextToken()
			continue
		}

		switch p.curToken.Type {
		case token.Default:
			if err := p.parseDefault(agent); err != nil {
				return nil, err
			}
		case token.Rules:
			if err := p.parseRules(agent); err != nil {
				return nil, err
			}
		default:
			p.skipUnknownStatementOrBlock()
		}
	}

	p.nextToken()
	return agent, nil
}

func (p *Parser) parseDefault(agent *Agent) error {
	p.nextToken()
	if !isEffectToken(p.curToken.Type) {
		return fmt.Errorf("expected effect after default at %d:%d", p.curToken.Line, p.curToken.Column)
	}
	agent.Default = p.curToken.Literal
	p.skipToStatementEnd()
	return nil
}

func (p *Parser) parseRules(agent *Agent) error {
	if !p.expectPeek(token.LBrace) {
		return p.unexpected("expected '{' after rules")
	}
	p.nextToken()

	for p.curToken.Type != token.RBrace {
		if p.curToken.Type == token.EOF {
			return fmt.Errorf("unterminated rules block")
		}
		if p.curToken.Type == token.Newline {
			p.nextToken()
			continue
		}
		if !isEffectToken(p.curToken.Type) {
			return fmt.Errorf("expected rule effect in rules block at %d:%d", p.curToken.Line, p.curToken.Column)
		}

		rule, err := p.parseRule()
		if err != nil {
			return err
		}
		agent.Rules = append(agent.Rules, *rule)
	}

	p.nextToken()
	return nil
}

func (p *Parser) parseRule() (*Rule, error) {
	rule := &Rule{Effect: p.curToken.Literal}
	p.nextToken()

	switch p.curToken.Type {
	case token.Ident, token.Star:
		rule.Tool = p.curToken.Literal
	default:
		return nil, fmt.Errorf("expected tool pattern after effect at %d:%d", p.curToken.Line, p.curToken.Column)
	}
	p.nextToken()

	for p.curToken.Type != token.Newline && p.curToken.Type != token.RBrace && p.curToken.Type != token.EOF {
		switch p.curToken.Type {
		case token.When:
			p.nextToken()
			expr, err := p.parseCondition()
			if err != nil {
				return nil, err
			}
			rule.Condition = expr
		case token.Notify:
			val, err := p.parseStringClause("notify")
			if err != nil {
				return nil, err
			}
			rule.Notify = val
		case token.Reason:
			val, err := p.parseStringClause("reason")
			if err != nil {
				return nil, err
			}
			rule.Reason = val
		case token.Reeval:
			val, err := p.parseBoolClause()
			if err != nil {
				return nil, err
			}
			rule.Reeval = &val
		default:
			return nil, fmt.Errorf("unexpected token %q in rule at %d:%d", p.curToken.Literal, p.curToken.Line, p.curToken.Column)
		}
	}

	if p.curToken.Type == token.Newline {
		p.nextToken()
	}
	return rule, nil
}

func (p *Parser) parseStringClause(name string) (string, error) {
	if !p.expectPeek(token.Colon) {
		return "", p.unexpected(fmt.Sprintf("expected ':' after %s", name))
	}
	if !p.expectPeek(token.String) {
		return "", p.unexpected(fmt.Sprintf("expected string after %s:", name))
	}
	val := p.curToken.Literal
	p.nextToken()
	return val, nil
}

func (p *Parser) parseBoolClause() (bool, error) {
	if !p.expectPeek(token.Colon) {
		return false, p.unexpected("expected ':' after reeval")
	}
	p.nextToken()
	if p.curToken.Type != token.True && p.curToken.Type != token.False {
		return false, p.unexpected("expected true|false after reeval:")
	}
	val := p.curToken.Type == token.True
	p.nextToken()
	return val, nil
}

func (p *Parser) parseCondition() (*Expr, error) {
	return p.parseOr()
}

func (p *Parser) parseOr() (*Expr, error) {
	left, err := p.parseAnd()
	if err != nil {
		return nil, err
	}
	for p.curToken.Type == token.Or {
		op := p.curToken.Literal
		p.nextToken()
		right, err := p.parseAnd()
		if err != nil {
			return nil, err
		}
		left = &Expr{Kind: "binary", Op: op, Left: left, Right: right}
	}
	return left, nil
}

func (p *Parser) parseAnd() (*Expr, error) {
	left, err := p.parseNot()
	if err != nil {
		return nil, err
	}
	for p.curToken.Type == token.And {
		op := p.curToken.Literal
		p.nextToken()
		right, err := p.parseNot()
		if err != nil {
			return nil, err
		}
		left = &Expr{Kind: "binary", Op: op, Left: left, Right: right}
	}
	return left, nil
}

func (p *Parser) parseNot() (*Expr, error) {
	if p.curToken.Type == token.Not {
		op := p.curToken.Literal
		p.nextToken()
		expr, err := p.parseNot()
		if err != nil {
			return nil, err
		}
		return &Expr{Kind: "unary", Op: op, Right: expr}, nil
	}
	return p.parseComparison()
}

func (p *Parser) parseComparison() (*Expr, error) {
	left, err := p.parsePrimary()
	if err != nil {
		return nil, err
	}

	if !isComparisonToken(p.curToken.Type) {
		return left, nil
	}

	op := p.curToken.Literal
	p.nextToken()
	right, err := p.parsePrimary()
	if err != nil {
		return nil, err
	}
	return &Expr{Kind: "binary", Op: op, Left: left, Right: right}, nil
}

func (p *Parser) parsePrimary() (*Expr, error) {
	switch p.curToken.Type {
	case token.Ident:
		e := &Expr{Kind: "ident", Value: p.curToken.Literal}
		p.nextToken()
		return e, nil
	case token.String, token.Number, token.Currency, token.True, token.False:
		e := &Expr{Kind: "literal", Value: p.curToken.Literal}
		p.nextToken()
		return e, nil
	case token.LParen:
		p.nextToken()
		expr, err := p.parseCondition()
		if err != nil {
			return nil, err
		}
		if p.curToken.Type != token.RParen {
			return nil, fmt.Errorf("expected ')' at %d:%d", p.curToken.Line, p.curToken.Column)
		}
		p.nextToken()
		return expr, nil
	case token.LBracket:
		return p.parseListLiteral()
	default:
		return nil, fmt.Errorf("unexpected token %q in condition at %d:%d", p.curToken.Literal, p.curToken.Line, p.curToken.Column)
	}
}

func (p *Parser) parseListLiteral() (*Expr, error) {
	items := ""
	for {
		p.nextToken()
		if p.curToken.Type == token.RBracket {
			break
		}
		if p.curToken.Type != token.String && p.curToken.Type != token.Number && p.curToken.Type != token.Currency && p.curToken.Type != token.True && p.curToken.Type != token.False && p.curToken.Type != token.Ident {
			return nil, fmt.Errorf("invalid list value %q at %d:%d", p.curToken.Literal, p.curToken.Line, p.curToken.Column)
		}
		if items != "" {
			items += ","
		}
		items += p.curToken.Literal
		p.nextToken()
		if p.curToken.Type == token.RBracket {
			break
		}
		if p.curToken.Type != token.Comma {
			return nil, fmt.Errorf("expected ',' in list literal at %d:%d", p.curToken.Line, p.curToken.Column)
		}
	}
	p.nextToken()
	return &Expr{Kind: "list", Value: items}, nil
}

func (p *Parser) skipUnknownStatementOrBlock() {
	if p.curToken.Type == token.LBrace {
		depth := 1
		for depth > 0 && p.curToken.Type != token.EOF {
			p.nextToken()
			if p.curToken.Type == token.LBrace {
				depth++
			}
			if p.curToken.Type == token.RBrace {
				depth--
			}
		}
		p.nextToken()
		return
	}
	p.skipToStatementEnd()
}

func (p *Parser) skipToStatementEnd() {
	for p.curToken.Type != token.Newline && p.curToken.Type != token.EOF && p.curToken.Type != token.RBrace {
		p.nextToken()
	}
	if p.curToken.Type == token.Newline {
		p.nextToken()
	}
}

func (p *Parser) consumeNewlines() {
	for p.curToken.Type == token.Newline {
		p.nextToken()
	}
}

func isEffectToken(t token.Type) bool {
	return t == token.Permit || t == token.Deny || t == token.DenyBang || t == token.Defer
}

func isComparisonToken(t token.Type) bool {
	switch t {
	case token.GT, token.GE, token.LT, token.LE, token.EQ, token.NE, token.Matches, token.In:
		return true
	default:
		return false
	}
}

func (p *Parser) expectPeek(t token.Type) bool {
	if p.peekToken.Type == t {
		p.nextToken()
		return true
	}
	return false
}

func (p *Parser) unexpected(msg string) error {
	return fmt.Errorf("%s at %d:%d", msg, p.peekToken.Line, p.peekToken.Column)
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}
