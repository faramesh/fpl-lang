package formatter

import (
	"sort"
	"strconv"
	"strings"
	"unicode"

	"github.com/faramesh/fpl-lang/reference/go/fpl/parser"
)

func FormatDocument(doc *parser.Document) string {
	var b strings.Builder

	for i, imp := range doc.Imports {
		if i > 0 {
			b.WriteString("\n")
		}
		b.WriteString("import ")
		b.WriteString(formatLiteral(imp.Ref))
		if imp.Alias != "" {
			b.WriteString(" as ")
			b.WriteString(formatLiteral(imp.Alias))
		}
		b.WriteString("\n")
	}

	if doc.Runtime != nil {
		writeConfigBlock(&b, "runtime", "", doc.Runtime.Fields)
	}

	for _, pb := range doc.Providers {
		writeConfigBlock(&b, "provider", pb.Name, pb.Fields)
	}

	for _, ib := range doc.Identities {
		writeConfigBlock(&b, "identity", ib.Name, ib.Fields)
	}

	if doc.Trust != nil && len(doc.Trust.Raw) > 0 {
		b.WriteString("trust {\n")
		for _, line := range doc.Trust.Raw {
			b.WriteString("  ")
			b.WriteString(strings.TrimSpace(line))
			b.WriteString("\n")
		}
		b.WriteString("}\n")
	}

	for i, ag := range doc.Agents {
		if i > 0 {
			b.WriteString("\n")
		}
		writeAgentBlock(&b, ag)
	}

	for _, sys := range doc.Systems {
		b.WriteString("\n")
		b.WriteString("system ")
		b.WriteString(formatLiteral(sys.ID))
		b.WriteString(" {\n")
		if sys.Version != "" {
			b.WriteString("  version ")
			b.WriteString(formatLiteral(sys.Version))
			b.WriteString("\n")
		}
		if sys.OnPolicyLoadFailure != "" {
			b.WriteString("  on_policy_load_failure ")
			b.WriteString(formatLiteral(sys.OnPolicyLoadFailure))
			b.WriteString("\n")
		}
		if sys.KillSwitchDefault != "" {
			b.WriteString("  kill_switch_default ")
			b.WriteString(formatLiteral(sys.KillSwitchDefault))
			b.WriteString("\n")
		}
		if sys.MaxOutputBytes != 0 {
			b.WriteString("  max_output_bytes ")
			b.WriteString(strconv.Itoa(sys.MaxOutputBytes))
			b.WriteString("\n")
		}
		b.WriteString("}\n")
	}

	for _, st := range doc.Topo {
		b.WriteString("manifest ")
		switch st.Kind {
		case parser.TopoOrchestrator:
			b.WriteString("orchestrator ")
			b.WriteString(formatLiteral(st.OrchID))
			b.WriteString(" undeclared ")
			b.WriteString(st.UndeclaredPolicy)
		case parser.TopoAllow:
			b.WriteString("grant ")
			b.WriteString(formatLiteral(st.AllowOrchID))
			b.WriteString(" to ")
			b.WriteString(formatLiteral(st.TargetAgentID))
			b.WriteString(" max ")
			b.WriteString(strconv.Itoa(st.MaxPerSession))
			if st.RequiresApproval {
				b.WriteString(" approval")
			}
		}
		b.WriteString("\n")
	}

	for _, rule := range doc.FlatRules {
		writeRule(&b, rule)
	}

	return strings.TrimRight(b.String(), "\n") + "\n"
}

func writeAgentBlock(b *strings.Builder, ag *parser.AgentBlock) {
	b.WriteString("agent ")
	b.WriteString(formatLiteral(ag.ID))
	b.WriteString(" {\n")
	if ag.Default != "" {
		b.WriteString("  default ")
		b.WriteString(ag.Default)
		b.WriteString("\n")
	}
	if ag.Model != "" {
		b.WriteString("  model ")
		b.WriteString(formatLiteral(ag.Model))
		b.WriteString("\n")
	}
	if ag.Framework != "" {
		b.WriteString("  framework ")
		b.WriteString(formatLiteral(ag.Framework))
		b.WriteString("\n")
	}
	if ag.Version != "" {
		b.WriteString("  version ")
		b.WriteString(formatLiteral(ag.Version))
		b.WriteString("\n")
	}

	if len(ag.Vars) > 0 {
		keys := sortedKeys(ag.Vars)
		for _, key := range keys {
			b.WriteString("  var ")
			b.WriteString(key)
			b.WriteString(" ")
			b.WriteString(formatLiteral(ag.Vars[key]))
			b.WriteString("\n")
		}
	}

	for _, bb := range ag.Budgets {
		b.WriteString("\n  budget ")
		b.WriteString(formatLiteral(bb.ID))
		b.WriteString(" {\n")
		if bb.Max != 0 {
			b.WriteString("    max $")
			b.WriteString(strconv.FormatFloat(bb.Max, 'f', -1, 64))
			b.WriteString("\n")
		}
		if bb.Daily != 0 {
			b.WriteString("    daily $")
			b.WriteString(strconv.FormatFloat(bb.Daily, 'f', -1, 64))
			b.WriteString("\n")
		}
		if bb.MaxCalls != 0 {
			b.WriteString("    max_calls ")
			b.WriteString(strconv.FormatInt(bb.MaxCalls, 10))
			b.WriteString("\n")
		}
		if bb.WarnAt != 0 {
			b.WriteString("    warn_at = ")
			b.WriteString(strconv.FormatFloat(bb.WarnAt, 'f', -1, 64))
			b.WriteString("\n")
		}
		if bb.OnExceed != "" {
			b.WriteString("    on_exceed ")
			b.WriteString(bb.OnExceed)
			b.WriteString("\n")
		}
		b.WriteString("  }\n")
	}

	for _, pb := range ag.Phases {
		b.WriteString("\n  phase ")
		b.WriteString(formatLiteral(pb.ID))
		b.WriteString(" {\n")
		for _, rule := range pb.Rules {
			writeRuleWithIndent(b, rule, "    ")
		}
		if pb.Duration != "" {
			b.WriteString("    duration ")
			b.WriteString(formatLiteral(pb.Duration))
			b.WriteString("\n")
		}
		if pb.Next != "" {
			b.WriteString("    next ")
			b.WriteString(formatLiteral(pb.Next))
			b.WriteString("\n")
		}
		b.WriteString("  }\n")
	}

	if len(ag.Rules) > 0 {
		b.WriteString("\n  rules {\n")
		for _, rule := range ag.Rules {
			writeRuleWithIndent(b, rule, "    ")
		}
		b.WriteString("  }\n")
	}

	for _, d := range ag.Delegates {
		b.WriteString("\n  delegate ")
		b.WriteString(formatLiteral(d.TargetAgent))
		b.WriteString(" {\n")
		if d.Scope != "" {
			b.WriteString("    scope ")
			b.WriteString(formatLiteral(d.Scope))
			b.WriteString("\n")
		}
		if d.TTL != "" {
			b.WriteString("    ttl ")
			b.WriteString(formatLiteral(d.TTL))
			b.WriteString("\n")
		}
		if d.Ceiling != "" {
			b.WriteString("    ceiling ")
			b.WriteString(formatLiteral(d.Ceiling))
			b.WriteString("\n")
		}
		b.WriteString("  }\n")
	}

	if len(ag.Ambients) > 0 {
		for _, amb := range ag.Ambients {
			b.WriteString("\n  ambient {\n")
			keys := sortedKeys(amb.Limits)
			for _, key := range keys {
				b.WriteString("    ")
				b.WriteString(key)
				b.WriteString(" ")
				b.WriteString(formatLiteral(amb.Limits[key]))
				b.WriteString("\n")
			}
			if amb.OnExceed != "" {
				b.WriteString("    on_exceed ")
				b.WriteString(formatLiteral(amb.OnExceed))
				b.WriteString("\n")
			}
			b.WriteString("  }\n")
		}
	}

	for _, sel := range ag.Selectors {
		b.WriteString("\n  selector ")
		b.WriteString(formatLiteral(sel.ID))
		b.WriteString(" {\n")
		if sel.Source != "" {
			b.WriteString("    source ")
			b.WriteString(formatLiteral(sel.Source))
			b.WriteString("\n")
		}
		if sel.Cache != "" {
			b.WriteString("    cache ")
			b.WriteString(formatLiteral(sel.Cache))
			b.WriteString("\n")
		}
		if sel.OnUnavailable != "" {
			b.WriteString("    on_unavailable ")
			b.WriteString(formatLiteral(sel.OnUnavailable))
			b.WriteString("\n")
		}
		if sel.OnTimeout != "" {
			b.WriteString("    on_timeout ")
			b.WriteString(formatLiteral(sel.OnTimeout))
			b.WriteString("\n")
		}
		b.WriteString("  }\n")
	}

	for _, cred := range ag.Credentials {
		b.WriteString("\n  credential ")
		b.WriteString(formatLiteral(cred.ID))
		b.WriteString(" {\n")
		if len(cred.Scope) > 0 {
			b.WriteString("    scope ")
			for i, s := range cred.Scope {
				if i > 0 {
					b.WriteByte(' ')
				}
				b.WriteString(formatLiteral(s))
			}
			b.WriteString("\n")
		}
		if cred.MaxScope != "" {
			b.WriteString("    max_scope ")
			b.WriteString(formatLiteral(cred.MaxScope))
			b.WriteString("\n")
		}
		if cred.Backend != "" {
			b.WriteString("    backend ")
			b.WriteString(formatLiteral(cred.Backend))
			b.WriteString("\n")
		}
		if cred.Path != "" {
			b.WriteString("    path ")
			b.WriteString(formatLiteral(cred.Path))
			b.WriteString("\n")
		}
		if cred.TTL != "" {
			b.WriteString("    ttl ")
			b.WriteString(formatLiteral(cred.TTL))
			b.WriteString("\n")
		}
		b.WriteString("  }\n")
	}

	if ag.Egress != nil {
		b.WriteString("\n  egress {\n")
		if len(ag.Egress.Allow) > 0 {
			b.WriteString("    allow = [")
			writeStringList(b, ag.Egress.Allow)
			b.WriteString("]\n")
		}
		if len(ag.Egress.Deny) > 0 {
			b.WriteString("    deny = [")
			writeStringList(b, ag.Egress.Deny)
			b.WriteString("]\n")
		}
		b.WriteString("  }\n")
	}

	if ag.ModelPolicy != nil && len(ag.ModelPolicy.Allow) > 0 {
		b.WriteString("\n  model_policy {\n")
		b.WriteString("    allow = [")
		writeStringList(b, ag.ModelPolicy.Allow)
		b.WriteString("]\n")
		b.WriteString("  }\n")
	}

	if ag.Session != nil && (ag.Session.MaxDuration != "" || ag.Session.IdleTimeout != "") {
		b.WriteString("\n  session {\n")
		if ag.Session.MaxDuration != "" {
			b.WriteString("    max_duration ")
			b.WriteString(formatLiteral(ag.Session.MaxDuration))
			b.WriteString("\n")
		}
		if ag.Session.IdleTimeout != "" {
			b.WriteString("    idle_timeout ")
			b.WriteString(formatLiteral(ag.Session.IdleTimeout))
			b.WriteString("\n")
		}
		b.WriteString("  }\n")
	}

	if ag.Spawn != nil && (ag.Spawn.MaxConcurrent != 0 || len(ag.Spawn.AllowedTypes) > 0) {
		b.WriteString("\n  spawn {\n")
		if ag.Spawn.MaxConcurrent != 0 {
			b.WriteString("    max_concurrent ")
			b.WriteString(strconv.Itoa(ag.Spawn.MaxConcurrent))
			b.WriteString("\n")
		}
		if len(ag.Spawn.AllowedTypes) > 0 {
			b.WriteString("    allowed_types = [")
			writeStringList(b, ag.Spawn.AllowedTypes)
			b.WriteString("]\n")
		}
		b.WriteString("  }\n")
	}

	if ag.CompletionGate != nil && len(ag.CompletionGate.Requires) > 0 {
		b.WriteString("\n  completion_gate {\n")
		for _, req := range ag.CompletionGate.Requires {
			b.WriteString("    require ")
			b.WriteString(formatLiteral(req))
			b.WriteString("\n")
		}
		b.WriteString("  }\n")
	}

	if ag.Enforcement != nil && len(ag.Enforcement.Fields) > 0 {
		b.WriteString("\n  enforcement {\n")
		keys := sortedKeysConfig(ag.Enforcement.Fields)
		for _, key := range keys {
			b.WriteString("    ")
			b.WriteString(key)
			b.WriteString(" ")
			b.WriteString(formatConfigValue(ag.Enforcement.Fields[key]))
			b.WriteString("\n")
		}
		b.WriteString("  }\n")
	}

	for _, al := range ag.Alerts {
		b.WriteString("\n  alert {\n")
		if al.On != "" {
			b.WriteString("    on ")
			b.WriteString(formatLiteral(al.On))
			b.WriteString("\n")
		}
		if al.Notify != "" {
			b.WriteString("    notify ")
			b.WriteString(formatLiteral(al.Notify))
			b.WriteString("\n")
		}
		b.WriteString("  }\n")
	}

	b.WriteString("}\n")
}

func writeRule(b *strings.Builder, rule *parser.Rule) {
	writeRuleWithIndent(b, rule, "")
}

func writeRuleWithIndent(b *strings.Builder, rule *parser.Rule, indent string) {
	b.WriteString(indent)
	b.WriteString(rule.Effect)
	b.WriteString(" ")
	b.WriteString(formatLiteral(rule.Tool))
	if rule.Condition != "" {
		b.WriteString(" when ")
		b.WriteString(rule.Condition)
	}
	if rule.Notify != "" {
		b.WriteString(" notify: ")
		b.WriteString(strconv.Quote(rule.Notify))
	}
	if rule.Reason != "" {
		b.WriteString(" reason: ")
		b.WriteString(strconv.Quote(rule.Reason))
	}
	if rule.Reeval != nil {
		b.WriteString(" reeval: ")
		if *rule.Reeval {
			b.WriteString("true")
		} else {
			b.WriteString("false")
		}
	}
	b.WriteString("\n")
}

func writeConfigBlock(b *strings.Builder, kind, name string, fields map[string]parser.ConfigValue) {
	b.WriteString(kind)
	if name != "" {
		b.WriteString(" ")
		b.WriteString(formatLiteral(name))
	}
	b.WriteString(" {\n")
	keys := sortedKeysConfig(fields)
	for _, key := range keys {
		b.WriteString("  ")
		b.WriteString(key)
		b.WriteString(" ")
		b.WriteString(formatConfigValue(fields[key]))
		b.WriteString("\n")
	}
	b.WriteString("}\n")
}

func formatConfigValue(v parser.ConfigValue) string {
	switch v.Kind {
	case parser.ConfigString, parser.ConfigIdent:
		return formatLiteral(v.String)
	case parser.ConfigNumber:
		return strconv.FormatFloat(v.Number, 'f', -1, 64)
	case parser.ConfigBool:
		if v.Bool {
			return "true"
		}
		return "false"
	case parser.ConfigEnv:
		return "env(" + formatLiteral(v.EnvVar) + ")"
	default:
		return ""
	}
}

func formatLiteral(v string) string {
	if needsQuotes(v) {
		return strconv.Quote(v)
	}
	return v
}

func needsQuotes(s string) bool {
	if strings.TrimSpace(s) == "" {
		return false
	}
	for _, r := range s {
		if unicode.IsSpace(r) || r == '"' {
			return true
		}
	}
	return !isBareFPLLiteral(s)
}

func isBareFPLLiteral(s string) bool {
	b := []byte(s)
	if len(b) == 0 {
		return false
	}
	if isIdentStart(b[0]) {
		for i := 1; i < len(b); i++ {
			if !isIdentCont(b[i]) {
				return false
			}
		}
		return true
	}
	if isDigit(b[0]) || (b[0] == '-' && len(b) > 1 && isDigit(b[1])) {
		i := 0
		if b[0] == '-' {
			i++
		}
		hasDot := false
		for i < len(b) && (isDigit(b[i]) || b[i] == '.') {
			if b[i] == '.' {
				if hasDot {
					return false
				}
				hasDot = true
			}
			i++
		}
		if i == len(b) {
			return true
		}
		return isUnitSuffix(string(b[i:]))
	}
	return false
}

func isDigit(c byte) bool { return c >= '0' && c <= '9' }

func isIdentStart(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_' ||
		c == '[' || c == ']' || c == ',' || c == '/' || c == '*' || c == '+' || c == '&' || c == '|'
}

func isIdentCont(c byte) bool {
	return isIdentStart(c) || isDigit(c) || c == '-' || c == '.' || c == '!' || c == '@' || c == '%' || c == '^'
}

func isUnitSuffix(s string) bool {
	switch strings.ToLower(s) {
	case "s", "ms", "m", "h", "d", "w",
		"b", "kb", "mb", "gb", "tb",
		"usd", "eur", "gbp":
		return true
	default:
		return false
	}
}

func writeStringList(b *strings.Builder, values []string) {
	for i, value := range values {
		if i > 0 {
			b.WriteByte(' ')
		}
		b.WriteString(formatLiteral(value))
	}
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func sortedKeysConfig(m map[string]parser.ConfigValue) []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
