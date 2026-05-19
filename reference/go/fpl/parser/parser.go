package parser

import (
	"fmt"
	"strings"
)

type parser struct {
	src    string
	tokens []token
	pos    int
}

type Parser struct {
	inner *parser
	err   error
}

type tokenKind int

const (
	tkIdent tokenKind = iota
	tkString
	tkNumber
	tkLBrace
	tkRBrace
	tkLParen
	tkRParen
	tkEq
	tkColon
	tkBang
	tkDollar
	tkLBracket
	tkRBracket
	tkEOF
)

type token struct {
	kind tokenKind
	val  string
	line int
}

func tokenize(src string) ([]token, error) {
	var tokens []token
	line := 1
	i := 0
	for i < len(src) {
		ch := src[i]
		switch {
		case ch == '\n':
			line++
			i++
		case ch == '\r':
			i++
		case ch == ' ' || ch == '\t':
			i++
		case ch == '#':
			for i < len(src) && src[i] != '\n' {
				i++
			}
		case ch == '{':
			tokens = append(tokens, token{tkLBrace, "{", line})
			i++
		case ch == '}':
			tokens = append(tokens, token{tkRBrace, "}", line})
			i++
		case ch == '(':
			tokens = append(tokens, token{tkLParen, "(", line})
			i++
		case ch == ')':
			tokens = append(tokens, token{tkRParen, ")", line})
			i++
		case ch == '=':
			if i+1 < len(src) && src[i+1] == '=' {
				tokens = append(tokens, token{tkIdent, "==", line})
				i += 2
			} else {
				tokens = append(tokens, token{tkEq, "=", line})
				i++
			}
		case ch == ':':
			tokens = append(tokens, token{tkColon, ":", line})
			i++
		case ch == '!' && i+1 < len(src) && src[i+1] == '=':
			tokens = append(tokens, token{tkIdent, "!=", line})
			i += 2
		case ch == '!' && (i+1 >= len(src) || src[i+1] == ' ' || src[i+1] == '\t' || src[i+1] == '\n' || src[i+1] == '\r'):
			tokens = append(tokens, token{tkBang, "!", line})
			i++
		case ch == '!':
			j := i
			for j < len(src) && isIdentCont(src[j]) {
				j++
			}
			tokens = append(tokens, token{tkIdent, src[i:j], line})
			i = j
		case ch == '<':
			if i+1 < len(src) && src[i+1] == '=' {
				tokens = append(tokens, token{tkIdent, "<=", line})
				i += 2
			} else {
				tokens = append(tokens, token{tkIdent, "<", line})
				i++
			}
		case ch == '>':
			if i+1 < len(src) && src[i+1] == '=' {
				tokens = append(tokens, token{tkIdent, ">=", line})
				i += 2
			} else {
				tokens = append(tokens, token{tkIdent, ">", line})
				i++
			}
		case ch == '$':
			tokens = append(tokens, token{tkDollar, "$", line})
			i++
		case ch == '[':
			tokens = append(tokens, token{tkLBracket, "[", line})
			i++
		case ch == ']':
			tokens = append(tokens, token{tkRBracket, "]", line})
			i++
		case ch == '"' || ch == '\'':
			quote := ch
			j := i + 1
			for j < len(src) && src[j] != quote {
				if src[j] == '\\' {
					j++
				}
				j++
			}
			if j >= len(src) {
				return nil, fmt.Errorf("line %d: unterminated string", line)
			}
			tokens = append(tokens, token{tkString, src[i+1 : j], line})
			i = j + 1
		case isDigit(ch) || (ch == '-' && i+1 < len(src) && isDigit(src[i+1])):
			j := i
			if ch == '-' {
				j++
			}
			for j < len(src) && (isDigit(src[j]) || src[j] == '.') {
				j++
			}
			tokens = append(tokens, token{tkNumber, src[i:j], line})
			i = j
		case isIdentStart(ch):
			j := i
			for j < len(src) && isIdentCont(src[j]) {
				j++
			}
			tokens = append(tokens, token{tkIdent, src[i:j], line})
			i = j
		default:
			return nil, fmt.Errorf("line %d: unexpected character %q", line, string(ch))
		}
	}
	tokens = append(tokens, token{tkEOF, "", line})
	return tokens, nil
}

func isDigit(c byte) bool { return c >= '0' && c <= '9' }

func isIdentStart(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_' ||
		c == '[' || c == ']' || c == ',' || c == '/' || c == '*' || c == '+' || c == '&' || c == '|'
}

func isIdentCont(c byte) bool {
	return isIdentStart(c) || isDigit(c) || c == '-' || c == '.' || c == '!' || c == '@' || c == '%' || c == '^'
}

func (p *parser) peek() token {
	if p.pos >= len(p.tokens) {
		return token{tkEOF, "", 0}
	}
	return p.tokens[p.pos]
}

func (p *parser) next() token {
	t := p.peek()
	if t.kind != tkEOF {
		p.pos++
	}
	return t
}

func (p *parser) expect(kind tokenKind) (token, error) {
	t := p.next()
	if t.kind != kind {
		return t, fmt.Errorf("line %d: expected %v, got %q", t.line, kind, t.val)
	}
	return t, nil
}

func (p *parser) expectIdent(val string) error {
	t := p.next()
	if t.kind != tkIdent || t.val != val {
		return fmt.Errorf("line %d: expected %q, got %q", t.line, val, t.val)
	}
	return nil
}

func (p *parser) peekIdent() string {
	t := p.peek()
	if t.kind == tkIdent {
		return t.val
	}
	return ""
}

func (p *parser) stringOrIdent() (string, error) {
	t := p.peek()
	switch t.kind {
	case tkString:
		p.next()
		return t.val, nil
	case tkIdent:
		p.next()
		return t.val, nil
	case tkNumber:
		p.next()
		if p.peek().kind == tkIdent && isUnitSuffix(p.peek().val) {
			unit := p.next()
			return t.val + unit.val, nil
		}
		return t.val, nil
	default:
		return "", fmt.Errorf("line %d: expected string or identifier, got %q", t.line, t.val)
	}
}

// ParseDocument parses a full structured FPL document.
func ParseDocument(src string) (*Document, error) {
	tokens, err := tokenize(src)
	if err != nil {
		return nil, err
	}
	p := &parser{src: src, tokens: tokens}
	return p.parseDocument()
}

func New(input string) *Parser {
	tokens, err := tokenize(input)
	if err != nil {
		return &Parser{err: err}
	}
	return &Parser{inner: &parser{src: input, tokens: tokens}}
}

func (p *Parser) ParseDocument() (*Document, error) {
	if p == nil || p.inner == nil {
		if p != nil && p.err != nil {
			return nil, p.err
		}
		return nil, fmt.Errorf("parser is not initialized")
	}
	if p.err != nil {
		return nil, p.err
	}
	return p.inner.parseDocument()
}

func (p *parser) parseDocument() (*Document, error) {
	doc := &Document{}

	for p.peek().kind != tkEOF {
		t := p.peek()
		if t.kind != tkIdent {
			return nil, fmt.Errorf("line %d: unexpected token %q", t.line, t.val)
		}

		switch t.val {
		case "import":
			imp, err := p.parseImportDecl()
			if err != nil {
				return nil, err
			}
			doc.Imports = append(doc.Imports, imp)
		case "runtime":
			rb, err := p.parseRuntimeBlock()
			if err != nil {
				return nil, err
			}
			doc.Runtime = rb
		case "provider":
			pb, err := p.parseNamedProviderBlock()
			if err != nil {
				return nil, err
			}
			doc.Providers = append(doc.Providers, pb)
		case "identity":
			ib, err := p.parseNamedIdentityBlock()
			if err != nil {
				return nil, err
			}
			doc.Identities = append(doc.Identities, ib)
		case "trust":
			tb, err := p.parseTrustBlock()
			if err != nil {
				return nil, err
			}
			doc.Trust = tb
		case "agent":
			ab, err := p.parseAgentBlock()
			if err != nil {
				return nil, err
			}
			doc.Agents = append(doc.Agents, ab)
		case "system":
			sb, err := p.parseSystemBlock()
			if err != nil {
				return nil, err
			}
			doc.Systems = append(doc.Systems, sb)
		case "manifest":
			topo, remaining, err := scanManifestLines(p.collectLine())
			if err != nil {
				return nil, err
			}
			doc.Topo = append(doc.Topo, topo...)
			_ = remaining
		case "permit", "allow", "approve", "deny", "block", "reject", "defer", "deny!":
			rule, err := p.parseFlatRule()
			if err != nil {
				return nil, err
			}
			doc.FlatRules = append(doc.FlatRules, rule)
		default:
			return nil, fmt.Errorf("line %d: unexpected keyword %q", t.line, t.val)
		}
	}

	if err := validateDocument(doc); err != nil {
		return nil, err
	}
	return doc, nil
}

func (p *parser) collectLine() string {
	start := p.pos
	line := p.peek().line
	for p.peek().kind != tkEOF && p.peek().line == line {
		p.next()
	}
	var parts []string
	for i := start; i < p.pos; i++ {
		parts = append(parts, p.tokens[i].val)
	}
	return strings.Join(parts, " ")
}

func (p *parser) parseFlatRule() (*Rule, error) {
	effectTok := p.next()
	effect := effectTok.val
	if effect == "deny" && p.peek().kind == tkBang {
		p.next()
		effect = "deny!"
	}

	tool, err := p.stringOrIdent()
	if err != nil {
		return nil, fmt.Errorf("line %d: rule tool: %w", effectTok.line, err)
	}

	rule := &Rule{Effect: effect, Tool: tool}

	for p.peek().kind != tkEOF {
		ident := p.peekIdent()
		switch ident {
		case "when":
			p.next()
			cond, err := p.consumeUntilKeyword()
			if err != nil {
				return nil, err
			}
			rule.Condition = cond
		case "notify":
			p.next()
			if p.peek().kind == tkColon {
				p.next()
			}
			v, err := p.stringOrIdent()
			if err != nil {
				return nil, err
			}
			rule.Notify = v
		case "reason":
			p.next()
			if p.peek().kind == tkColon {
				p.next()
			}
			v, err := p.stringOrIdent()
			if err != nil {
				return nil, err
			}
			rule.Reason = v
		case "reeval":
			p.next()
			val, err := p.parseBoolClause()
			if err != nil {
				return nil, err
			}
			rule.Reeval = &val
		case "host":
			p.next()
			v, err := p.consumeFlatRuleClauseValue()
			if err != nil {
				return nil, err
			}
			rule.Host = strings.TrimSpace(v)
		case "port":
			p.next()
			v, err := p.consumeFlatRuleClauseValue()
			if err != nil {
				return nil, err
			}
			rule.Port = strings.TrimSpace(v)
		case "method":
			p.next()
			v, err := p.consumeFlatRuleClauseValue()
			if err != nil {
				return nil, err
			}
			rule.Method = strings.TrimSpace(v)
		case "path":
			p.next()
			v, err := p.consumeFlatRuleClauseValue()
			if err != nil {
				return nil, err
			}
			rule.Path = strings.TrimSpace(v)
		case "query":
			p.next()
			v, err := p.consumeFlatRuleClauseValue()
			if err != nil {
				return nil, err
			}
			k, qv, err := parseFlatRuleKeyValue(v)
			if err != nil {
				return nil, fmt.Errorf("line %d: query clause: %w", p.peek().line, err)
			}
			if rule.Query == nil {
				rule.Query = make(map[string]string)
			}
			rule.Query[k] = qv
		case "header", "headers":
			p.next()
			v, err := p.consumeFlatRuleClauseValue()
			if err != nil {
				return nil, err
			}
			k, hv, err := parseFlatRuleKeyValue(v)
			if err != nil {
				return nil, fmt.Errorf("line %d: header clause: %w", p.peek().line, err)
			}
			if rule.Headers == nil {
				rule.Headers = make(map[string]string)
			}
			rule.Headers[k] = hv
		default:
			return rule, nil
		}
	}
	return rule, nil
}
func (p *parser) consumeFlatRuleClauseValue() (string, error) {
	if p.peek().kind == tkColon {
		p.next()
	}
	return p.consumeUntilKeyword()
}

func parseFlatRuleKeyValue(raw string) (string, string, error) {
	parts := strings.SplitN(strings.TrimSpace(raw), "=", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("expected key=value")
	}
	key := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])
	if key == "" {
		return "", "", fmt.Errorf("empty key")
	}
	if value == "" {
		return "", "", fmt.Errorf("empty value")
	}
	return key, value, nil
}

func (p *parser) consumeUntilKeyword() (string, error) {
	startLine := p.peek().line
	var parts []string
	for {
		t := p.peek()
		if t.kind == tkEOF || t.kind == tkLBrace || t.kind == tkRBrace {
			break
		}
		if t.kind == tkIdent && isFlatRuleClauseKeyword(t.val) {
			break
		}
		if t.kind == tkIdent && isEffectKeyword(t.val) && t.line != startLine {
			break
		}
		if t.kind == tkIdent && isTopLevelKeyword(t.val) {
			break
		}
		p.next()
		val := t.val
		switch t.kind {
		case tkString:
			val = `"` + val + `"`
		case tkColon:
			val = ":"
		case tkDollar:
			val = "$"
		case tkBang:
			val = "!"
		case tkLParen:
			val = "("
		case tkRParen:
			val = ")"
		case tkEq:
			val = "="
		}
		parts = append(parts, val)
	}
	return formatExprParts(parts), nil
}

func (p *parser) parseBoolClause() (bool, error) {
	if p.peek().kind == tkColon {
		p.next()
	}
	t := p.next()
	if t.kind != tkIdent || (t.val != "true" && t.val != "false") {
		return false, fmt.Errorf("line %d: expected true|false after reeval:", t.line)
	}
	return t.val == "true", nil
}

func formatExprParts(parts []string) string {
	var b strings.Builder
	for i, part := range parts {
		if i == 0 {
			b.WriteString(part)
			continue
		}
		prev := parts[i-1]
		switch {
		case part == "(":
			b.WriteString(part)
		case part == ")" || part == ",":
			b.WriteString(part)
		case part == "=" && prev == "=":
			continue
		case prev == "(":
			b.WriteString(part)
		default:
			b.WriteByte(' ')
			b.WriteString(part)
		}
	}
	return b.String()
}

func isFlatRuleClauseKeyword(v string) bool {
	switch v {
	case "when", "notify", "reason", "host", "port", "method", "path", "query", "header", "headers", "reeval":
		return true
	default:
		return false
	}
}

func isTopLevelKeyword(s string) bool {
	switch s {
	case "import", "runtime", "provider", "identity", "trust",
		"agent", "system", "permit", "allow", "approve", "deny", "deny!", "block", "reject", "defer", "manifest":
		return true
	}
	return false
}

func isCredentialKeyword(s string) bool {
	switch s {
	case "scope", "max_scope", "backend", "path", "ttl":
		return true
	default:
		return false
	}
}

func isUnitSuffix(s string) bool {
	switch strings.ToLower(s) {
	case "s", "ms", "m", "h", "d", "w",
		"b", "kb", "mb", "gb", "tb",
		"usd", "eur", "gbp":
		return true
	}
	return false
}

func isEffectKeyword(s string) bool {
	switch s {
	case "permit", "allow", "approve", "deny", "deny!", "block", "reject", "defer":
		return true
	default:
		return false
	}
}

func validateDocument(doc *Document) error {
	for _, agent := range doc.Agents {
		if err := validateAgent(agent); err != nil {
			return err
		}
	}
	return nil
}

func validateAgent(agent *AgentBlock) error {
	denyBangTools := make(map[string]int)

	for i, rule := range agent.Rules {
		switch rule.Effect {
		case "deny!":
			denyBangTools[rule.Tool] = i
		case "permit":
			if denyIndex, ok := denyBangTools[rule.Tool]; ok {
				return fmt.Errorf(
					"policy violates deny! invariant in agent %q: permit rule for tool pattern %q follows deny! rule at rule index %d",
					agent.ID,
					rule.Tool,
					denyIndex+1,
				)
			}
		}
	}

	return nil
}
