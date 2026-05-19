package parser

// Rule is a single flat rule line.
type Rule struct {
	Effect    string
	Tool      string
	Condition string
	Notify    string
	Reason    string
	Host      string
	Port      string
	Method    string
	Path      string
	Query     map[string]string
	Headers   map[string]string
	Reeval    *bool
}

// Document is the top-level AST node for a structured FPL file.
type Document struct {
	Imports    []ImportDecl
	Runtime    *RuntimeBlock
	Providers  []*NamedProviderBlock
	Identities []*NamedIdentityBlock
	Trust      *TrustBlock
	Agents     []*AgentBlock
	Systems    []*SystemBlock
	FlatRules  []*Rule
	Topo       []TopoStatement
}

type AgentBlock struct {
	ID          string
	Default     string
	Model       string
	Framework   string
	Version     string
	Budgets     []*BudgetBlock
	Phases      []*PhaseBlock
	Rules       []*Rule
	Delegates   []*DelegateBlock
	Ambients    []*AmbientBlock
	Selectors   []*SelectorBlock
	Credentials []*CredentialBlock
	Vars        map[string]string

	RateLimits     []*RateLimitLine
	Redactions     []*RedactLine
	Egress         *EgressBlock
	ModelPolicy    *ModelPolicyBlock
	Session        *SessionBlock
	Spawn          *SpawnBlock
	CompletionGate *CompletionGateBlock
	Enforcement    *EnforcementBlock
	Alerts         []*AlertBlock
}

type SystemBlock struct {
	ID                  string
	Version             string
	OnPolicyLoadFailure string
	KillSwitchDefault   string
	MaxOutputBytes      int
}

type BudgetBlock struct {
	ID       string
	Max      float64
	Daily    float64
	MaxCalls int64
	WarnAt   float64
	OnExceed string
}

type PhaseBlock struct {
	ID       string
	Tools    []string
	Rules    []*Rule
	Duration string
	Next     string
}

type DelegateBlock struct {
	TargetAgent string
	Scope       string
	TTL         string
	Ceiling     string
}

type AmbientBlock struct {
	Limits   map[string]string
	OnExceed string
}

type SelectorBlock struct {
	ID            string
	Source        string
	Cache         string
	OnUnavailable string
	OnTimeout     string
}

type CredentialBlock struct {
	ID       string
	Scope    []string
	MaxScope string
	Backend  string
	Path     string
	TTL      string
}

// ConfigValue is a scalar in runtime, provider, or identity blocks.
type ConfigValue struct {
	Kind ConfigValueKind

	String string
	Number float64
	Bool   bool
	EnvVar string
}

type ConfigValueKind string

const (
	ConfigString ConfigValueKind = "string"
	ConfigNumber ConfigValueKind = "number"
	ConfigBool   ConfigValueKind = "bool"
	ConfigEnv    ConfigValueKind = "env"
	ConfigIdent  ConfigValueKind = "ident"
)

type ImportDecl struct {
	Ref   string
	Alias string
	Line  int
}

type RuntimeBlock struct {
	Fields map[string]ConfigValue
}

type NamedProviderBlock struct {
	Name   string
	Fields map[string]ConfigValue
}

type NamedIdentityBlock struct {
	Name   string
	Fields map[string]ConfigValue
}

type TrustBlock struct {
	Raw []string
}

type RateLimitLine struct {
	Pattern string
	Limit   int64
	Window  string
	Line    int
}

type RedactLine struct {
	Tool  string
	Paths []string
	Line  int
}

type EgressBlock struct {
	Allow []string
	Deny  []string
}

type ModelPolicyBlock struct {
	Allow []string
}

type SessionBlock struct {
	MaxDuration string
	IdleTimeout string
}

type SpawnBlock struct {
	MaxConcurrent int
	AllowedTypes  []string
}

type CompletionGateBlock struct {
	Requires []string
}

type AlertBlock struct {
	On     string
	Notify string
}

type EnforcementBlock struct {
	Fields map[string]ConfigValue
}

// TopoKind classifies topology statements parsed from FPL manifest lines.
type TopoKind int

const (
	TopoOrchestrator TopoKind = iota
	TopoAllow
)

// TopoStatement is a single manifest line after parse (decl or allow).
type TopoStatement struct {
	Kind TopoKind

	OrchID           string
	UndeclaredPolicy string

	AllowOrchID      string
	TargetAgentID    string
	MaxPerSession    int
	RequiresApproval bool
}
