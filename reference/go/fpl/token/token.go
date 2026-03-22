package token

type Type string

const (
	Illegal Type = "ILLEGAL"
	EOF     Type = "EOF"

	Ident    Type = "IDENT"
	String   Type = "STRING"
	Number   Type = "NUMBER"
	Currency Type = "CURRENCY"

	LBrace   Type = "LBRACE"
	RBrace   Type = "RBRACE"
	LParen   Type = "LPAREN"
	RParen   Type = "RPAREN"
	LBracket Type = "LBRACKET"
	RBracket Type = "RBRACKET"
	Comma    Type = "COMMA"
	Colon    Type = "COLON"
	Star     Type = "STAR"
	Newline  Type = "NEWLINE"

	GT Type = "GT"
	GE Type = "GE"
	LT Type = "LT"
	LE Type = "LE"
	EQ Type = "EQ"
	NE Type = "NE"

	Agent   Type = "AGENT"
	Default Type = "DEFAULT"
	Rules   Type = "RULES"

	Permit   Type = "PERMIT"
	Deny     Type = "DENY"
	DenyBang Type = "DENY_BANG"
	Defer    Type = "DEFER"

	When   Type = "WHEN"
	Notify Type = "NOTIFY"
	Reason Type = "REASON"
	Reeval Type = "REEVAL"

	And     Type = "AND"
	Or      Type = "OR"
	Not     Type = "NOT"
	In      Type = "IN"
	Matches Type = "MATCHES"

	True  Type = "TRUE"
	False Type = "FALSE"
)

type Token struct {
	Type    Type
	Literal string
	Line    int
	Column  int
}

var keywords = map[string]Type{
	"agent":   Agent,
	"default": Default,
	"rules":   Rules,

	"permit": Permit,
	"deny":   Deny,
	"defer":  Defer,

	"when":   When,
	"notify": Notify,
	"reason": Reason,
	"reeval": Reeval,

	"and":     And,
	"or":      Or,
	"not":     Not,
	"in":      In,
	"matches": Matches,

	"true":  True,
	"false": False,
}

func LookupIdent(ident string) Type {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return Ident
}
