package parser

import (
	"fmt"
	"strconv"
)

func (p *parser) parseAgentBlock() (*AgentBlock, error) {
	if err := p.expectIdent("agent"); err != nil {
		return nil, err
	}
	id, err := p.stringOrIdent()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(tkLBrace); err != nil {
		return nil, err
	}

	ab := &AgentBlock{ID: id, Vars: make(map[string]string)}

	for p.peek().kind != tkRBrace && p.peek().kind != tkEOF {
		kw := p.peekIdent()
		switch kw {
		case "default":
			p.next()
			v, err := p.stringOrIdent()
			if err != nil {
				return nil, err
			}
			ab.Default = v
		case "model":
			p.next()
			v, err := p.stringOrIdent()
			if err != nil {
				return nil, err
			}
			ab.Model = v
		case "framework":
			p.next()
			v, err := p.stringOrIdent()
			if err != nil {
				return nil, err
			}
			ab.Framework = v
		case "version":
			p.next()
			v, err := p.stringOrIdent()
			if err != nil {
				return nil, err
			}
			ab.Version = v
		case "var":
			p.next()
			name, err := p.stringOrIdent()
			if err != nil {
				return nil, err
			}
			val, err := p.parseVarValue()
			if err != nil {
				return nil, err
			}
			ab.Vars[name] = val
		case "budget":
			bb, err := p.parseBudgetBlock()
			if err != nil {
				return nil, err
			}
			ab.Budgets = append(ab.Budgets, bb)
		case "phase":
			pb, err := p.parsePhaseBlock()
			if err != nil {
				return nil, err
			}
			ab.Phases = append(ab.Phases, pb)
		case "rules":
			rules, err := p.parseRulesBlock()
			if err != nil {
				return nil, err
			}
			ab.Rules = append(ab.Rules, rules...)
		case "delegate":
			db, err := p.parseDelegateBlock()
			if err != nil {
				return nil, err
			}
			ab.Delegates = append(ab.Delegates, db)
		case "ambient":
			amb, err := p.parseAmbientBlock()
			if err != nil {
				return nil, err
			}
			ab.Ambients = append(ab.Ambients, amb)
		case "selector":
			sel, err := p.parseSelectorBlock()
			if err != nil {
				return nil, err
			}
			ab.Selectors = append(ab.Selectors, sel)
		case "credential":
			cred, err := p.parseCredentialBlock()
			if err != nil {
				return nil, err
			}
			ab.Credentials = append(ab.Credentials, cred)
		case "rate_limit":
			rl, err := p.parseRateLimitLine()
			if err != nil {
				return nil, err
			}
			ab.RateLimits = append(ab.RateLimits, rl)
		case "redact":
			rd, err := p.parseRedactLine()
			if err != nil {
				return nil, err
			}
			ab.Redactions = append(ab.Redactions, rd)
		case "egress":
			eg, err := p.parseEgressBlock()
			if err != nil {
				return nil, err
			}
			ab.Egress = eg
		case "model_policy":
			mp, err := p.parseModelPolicyBlock()
			if err != nil {
				return nil, err
			}
			ab.ModelPolicy = mp
		case "session":
			sb, err := p.parseSessionBlock()
			if err != nil {
				return nil, err
			}
			ab.Session = sb
		case "spawn":
			sp, err := p.parseSpawnBlock()
			if err != nil {
				return nil, err
			}
			ab.Spawn = sp
		case "completion_gate":
			cg, err := p.parseCompletionGateBlock()
			if err != nil {
				return nil, err
			}
			ab.CompletionGate = cg
		case "enforcement":
			enf, err := p.parseEnforcementBlock()
			if err != nil {
				return nil, err
			}
			ab.Enforcement = enf
		case "alert":
			al, err := p.parseAlertBlock()
			if err != nil {
				return nil, err
			}
			ab.Alerts = append(ab.Alerts, al)
		case "permit", "allow", "approve", "deny", "block", "reject", "defer", "deny!":
			rule, err := p.parseFlatRule()
			if err != nil {
				return nil, err
			}
			ab.Rules = append(ab.Rules, rule)
		default:
			return nil, fmt.Errorf("line %d: unexpected keyword %q in agent block", p.peek().line, kw)
		}
	}

	if _, err := p.expect(tkRBrace); err != nil {
		return nil, fmt.Errorf("agent %s: %w", id, err)
	}
	return ab, nil
}

func (p *parser) parseSystemBlock() (*SystemBlock, error) {
	if err := p.expectIdent("system"); err != nil {
		return nil, err
	}
	id, err := p.stringOrIdent()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(tkLBrace); err != nil {
		return nil, err
	}

	sb := &SystemBlock{ID: id}

	for p.peek().kind != tkRBrace && p.peek().kind != tkEOF {
		kw := p.peekIdent()
		switch kw {
		case "version":
			p.next()
			v, err := p.stringOrIdent()
			if err != nil {
				return nil, err
			}
			sb.Version = v
		case "on_policy_load_failure":
			p.next()
			v, err := p.stringOrIdent()
			if err != nil {
				return nil, err
			}
			sb.OnPolicyLoadFailure = v
		case "kill_switch_default":
			p.next()
			v, err := p.stringOrIdent()
			if err != nil {
				return nil, err
			}
			sb.KillSwitchDefault = v
		case "max_output_bytes":
			p.next()
			t := p.next()
			n, err := strconv.Atoi(t.val)
			if err != nil {
				return nil, fmt.Errorf("line %d: max_output_bytes: %w", t.line, err)
			}
			sb.MaxOutputBytes = n
		default:
			return nil, fmt.Errorf("line %d: unexpected keyword %q in system block", p.peek().line, kw)
		}
	}

	if _, err := p.expect(tkRBrace); err != nil {
		return nil, fmt.Errorf("system %s: %w", id, err)
	}
	return sb, nil
}

func (p *parser) parseBudgetBlock() (*BudgetBlock, error) {
	if err := p.expectIdent("budget"); err != nil {
		return nil, err
	}
	id, err := p.stringOrIdent()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(tkLBrace); err != nil {
		return nil, err
	}

	bb := &BudgetBlock{ID: id}

	for p.peek().kind != tkRBrace && p.peek().kind != tkEOF {
		kw := p.peekIdent()
		switch kw {
		case "max":
			p.next()
			v, err := p.parseCurrency()
			if err != nil {
				return nil, err
			}
			bb.Max = v
		case "daily":
			p.next()
			v, err := p.parseCurrency()
			if err != nil {
				return nil, err
			}
			bb.Daily = v
		case "max_calls":
			p.next()
			t := p.next()
			n, err := strconv.ParseInt(t.val, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("line %d: max_calls: %w", t.line, err)
			}
			bb.MaxCalls = n
		case "on_exceed":
			p.next()
			v, err := p.stringOrIdent()
			if err != nil {
				return nil, err
			}
			bb.OnExceed = v
		case "warn_at":
			p.next()
			p.skipOptionalEquals()
			t := p.next()
			v, err := strconv.ParseFloat(t.val, 64)
			if err != nil {
				return nil, fmt.Errorf("line %d: warn_at: %w", t.line, err)
			}
			bb.WarnAt = v
		default:
			return nil, fmt.Errorf("line %d: unexpected keyword %q in budget block", p.peek().line, kw)
		}
	}

	if _, err := p.expect(tkRBrace); err != nil {
		return nil, fmt.Errorf("budget %s: %w", id, err)
	}
	return bb, nil
}

func (p *parser) parsePhaseBlock() (*PhaseBlock, error) {
	if err := p.expectIdent("phase"); err != nil {
		return nil, err
	}
	id, err := p.stringOrIdent()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(tkLBrace); err != nil {
		return nil, err
	}

	pb := &PhaseBlock{ID: id}

	for p.peek().kind != tkRBrace && p.peek().kind != tkEOF {
		kw := p.peekIdent()
		switch kw {
		case "permit", "allow", "approve", "deny", "deny!", "block", "reject", "defer":
			rule, err := p.parseFlatRule()
			if err != nil {
				return nil, err
			}
			pb.Rules = append(pb.Rules, rule)
			pb.Tools = append(pb.Tools, rule.Tool)
		case "duration":
			p.next()
			v, err := p.stringOrIdent()
			if err != nil {
				return nil, err
			}
			pb.Duration = v
		case "next":
			p.next()
			v, err := p.stringOrIdent()
			if err != nil {
				return nil, err
			}
			pb.Next = v
		default:
			return nil, fmt.Errorf("line %d: unexpected keyword %q in phase block", p.peek().line, kw)
		}
	}

	if _, err := p.expect(tkRBrace); err != nil {
		return nil, fmt.Errorf("phase %s: %w", id, err)
	}
	return pb, nil
}

func (p *parser) parseRulesBlock() ([]*Rule, error) {
	if err := p.expectIdent("rules"); err != nil {
		return nil, err
	}
	if _, err := p.expect(tkLBrace); err != nil {
		return nil, err
	}

	var rules []*Rule
	for p.peek().kind != tkRBrace && p.peek().kind != tkEOF {
		rule, err := p.parseFlatRule()
		if err != nil {
			return nil, err
		}
		rules = append(rules, rule)
	}

	if _, err := p.expect(tkRBrace); err != nil {
		return nil, err
	}
	return rules, nil
}

func (p *parser) parseDelegateBlock() (*DelegateBlock, error) {
	if err := p.expectIdent("delegate"); err != nil {
		return nil, err
	}
	target, err := p.stringOrIdent()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(tkLBrace); err != nil {
		return nil, err
	}

	db := &DelegateBlock{TargetAgent: target}

	for p.peek().kind != tkRBrace && p.peek().kind != tkEOF {
		kw := p.peekIdent()
		switch kw {
		case "scope":
			p.next()
			v, err := p.stringOrIdent()
			if err != nil {
				return nil, err
			}
			db.Scope = v
		case "ttl":
			p.next()
			v, err := p.stringOrIdent()
			if err != nil {
				return nil, err
			}
			db.TTL = v
		case "ceiling":
			p.next()
			v, err := p.stringOrIdent()
			if err != nil {
				return nil, err
			}
			db.Ceiling = v
		default:
			return nil, fmt.Errorf("line %d: unexpected keyword %q in delegate block", p.peek().line, kw)
		}
	}

	if _, err := p.expect(tkRBrace); err != nil {
		return nil, fmt.Errorf("delegate %s: %w", target, err)
	}
	return db, nil
}

func (p *parser) parseAmbientBlock() (*AmbientBlock, error) {
	if err := p.expectIdent("ambient"); err != nil {
		return nil, err
	}
	if _, err := p.expect(tkLBrace); err != nil {
		return nil, err
	}

	ab := &AmbientBlock{Limits: make(map[string]string)}

	for p.peek().kind != tkRBrace && p.peek().kind != tkEOF {
		kw := p.peekIdent()
		if kw == "on_exceed" {
			p.next()
			v, err := p.stringOrIdent()
			if err != nil {
				return nil, err
			}
			ab.OnExceed = v
		} else if kw != "" {
			p.next()
			v, err := p.stringOrIdent()
			if err != nil {
				return nil, err
			}
			ab.Limits[kw] = v
		} else {
			return nil, fmt.Errorf("line %d: unexpected token in ambient block", p.peek().line)
		}
	}

	if _, err := p.expect(tkRBrace); err != nil {
		return nil, err
	}
	return ab, nil
}

func (p *parser) parseSelectorBlock() (*SelectorBlock, error) {
	if err := p.expectIdent("selector"); err != nil {
		return nil, err
	}
	id, err := p.stringOrIdent()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(tkLBrace); err != nil {
		return nil, err
	}

	sb := &SelectorBlock{ID: id}

	for p.peek().kind != tkRBrace && p.peek().kind != tkEOF {
		kw := p.peekIdent()
		switch kw {
		case "source":
			p.next()
			v, err := p.stringOrIdent()
			if err != nil {
				return nil, err
			}
			sb.Source = v
		case "cache":
			p.next()
			v, err := p.stringOrIdent()
			if err != nil {
				return nil, err
			}
			sb.Cache = v
		case "on_unavailable":
			p.next()
			v, err := p.stringOrIdent()
			if err != nil {
				return nil, err
			}
			sb.OnUnavailable = v
		case "on_timeout":
			p.next()
			v, err := p.stringOrIdent()
			if err != nil {
				return nil, err
			}
			sb.OnTimeout = v
		default:
			return nil, fmt.Errorf("line %d: unexpected keyword %q in selector block", p.peek().line, kw)
		}
	}

	if _, err := p.expect(tkRBrace); err != nil {
		return nil, fmt.Errorf("selector %s: %w", id, err)
	}
	return sb, nil
}

func (p *parser) parseCredentialBlock() (*CredentialBlock, error) {
	if err := p.expectIdent("credential"); err != nil {
		return nil, err
	}
	id, err := p.stringOrIdent()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(tkLBrace); err != nil {
		return nil, err
	}

	cb := &CredentialBlock{ID: id}

	for p.peek().kind != tkRBrace && p.peek().kind != tkEOF {
		kw := p.peekIdent()
		switch kw {
		case "scope":
			p.next()
			for p.peek().kind == tkIdent || p.peek().kind == tkString {
				if p.peek().kind == tkIdent && isCredentialKeyword(p.peek().val) {
					break
				}
				v, _ := p.stringOrIdent()
				cb.Scope = append(cb.Scope, v)
			}
		case "max_scope":
			p.next()
			v, err := p.stringOrIdent()
			if err != nil {
				return nil, err
			}
			cb.MaxScope = v
		case "backend":
			p.next()
			v, err := p.stringOrIdent()
			if err != nil {
				return nil, err
			}
			cb.Backend = v
		case "path":
			p.next()
			v, err := p.stringOrIdent()
			if err != nil {
				return nil, err
			}
			cb.Path = v
		case "ttl":
			p.next()
			v, err := p.stringOrIdent()
			if err != nil {
				return nil, err
			}
			cb.TTL = v
		default:
			return nil, fmt.Errorf("line %d: unexpected keyword %q in credential block", p.peek().line, kw)
		}
	}

	if _, err := p.expect(tkRBrace); err != nil {
		return nil, fmt.Errorf("credential %s: %w", id, err)
	}
	return cb, nil
}
