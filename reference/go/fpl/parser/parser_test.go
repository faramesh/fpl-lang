package parser

import "testing"

func TestParseDocumentOK(t *testing.T) {
	src := `agent payment-bot {
  default deny
  model "gpt-4o"
  framework "langgraph"

  budget session {
    max $500
    daily $2000
    max_calls 100
    on_exceed deny
  }

  phase intake {
    permit read_customer
    permit get_order
  }

  rules {
		deny! shell/run reason: "dangerous"
		permit read_customer when amount <= 500 and not cmd matches "rm -rf" notify: "finance" reeval: true
  }
}`

	p := New(src)
	doc, err := p.ParseDocument()
	if err != nil {
		t.Fatalf("ParseDocument error: %v", err)
	}
	if len(doc.Agents) != 1 {
		t.Fatalf("expected one agent, got %d", len(doc.Agents))
	}
	if doc.Agents[0].ID != "payment-bot" {
		t.Fatalf("unexpected agent id: %q", doc.Agents[0].ID)
	}
	ag := doc.Agents[0]
	if ag.Default != "deny" {
		t.Fatalf("unexpected default: %q", ag.Default)
	}
	if ag.Model != "gpt-4o" || ag.Framework != "langgraph" {
		t.Fatalf("unexpected agent metadata: %+v", ag)
	}
	if len(ag.Budgets) != 1 || ag.Budgets[0].ID != "session" || ag.Budgets[0].Max != 500 || ag.Budgets[0].Daily != 2000 || ag.Budgets[0].MaxCalls != 100 {
		t.Fatalf("unexpected budget block: %+v", ag.Budgets)
	}
	if len(ag.Phases) != 1 || ag.Phases[0].ID != "intake" || len(ag.Phases[0].Rules) != 2 {
		t.Fatalf("unexpected phase block: %+v", ag.Phases)
	}
	if len(ag.Rules) != 2 {
		t.Fatalf("expected two rules, got %d", len(ag.Rules))
	}
	if ag.Rules[0].Effect != "deny!" || ag.Rules[0].Tool != "shell/run" {
		t.Fatalf("unexpected first rule: %+v", ag.Rules[0])
	}
	if ag.Rules[1].Condition == "" {
		t.Fatal("expected condition on second rule")
	}
	if ag.Rules[1].Notify != "finance" {
		t.Fatalf("unexpected notify: %q", ag.Rules[1].Notify)
	}
	if ag.Rules[1].Reeval == nil || !*ag.Rules[1].Reeval {
		t.Fatalf("expected reeval true, got: %+v", ag.Rules[1].Reeval)
	}
}

func TestParseDocumentRejectsMissingAgentName(t *testing.T) {
	src := `agent {
  default deny
}`

	p := New(src)
	_, err := p.ParseDocument()
	if err == nil {
		t.Fatal("expected parse error for missing agent name")
	}
}

func TestParseDocumentRejectsBadRule(t *testing.T) {
	src := `agent payment-bot {
  rules {
    permit
  }
}`

	p := New(src)
	_, err := p.ParseDocument()
	if err == nil {
		t.Fatal("expected parse error for incomplete rule")
	}
}

func TestParseDocumentCanonicalizesEffectAliases(t *testing.T) {
	src := `agent payment-bot {
  default allow
  rules {
    approve read_customer
    block shell/run
  }
}`

	p := New(src)
	doc, err := p.ParseDocument()
	if err != nil {
		t.Fatalf("ParseDocument error: %v", err)
	}
	ag := doc.Agents[0]
	if ag.Default != "allow" {
		t.Fatalf("expected default allow to be preserved, got %q", ag.Default)
	}
	if len(ag.Rules) != 2 {
		t.Fatalf("expected two rules, got %d", len(ag.Rules))
	}
	if ag.Rules[0].Effect != "approve" {
		t.Fatalf("expected approve to be preserved, got %q", ag.Rules[0].Effect)
	}
	if ag.Rules[1].Effect != "block" {
		t.Fatalf("expected block to be preserved, got %q", ag.Rules[1].Effect)
	}
}

func TestParseDocumentRejectsDenyBangOverride(t *testing.T) {
	src := `agent payment-bot {
  rules {
    deny! shell/run
    permit shell/run
  }
}`

	p := New(src)
	_, err := p.ParseDocument()
	if err == nil {
		t.Fatal("expected parse error for deny! override")
	}
}

func TestParseDocumentSkipsStructuredBlocks(t *testing.T) {
	src := `agent support-bot {
	default deny

	model "gpt-4o"

	budget session {
		max $200
		max_calls 50
		on_exceed deny
	}

	phase intake {
		permit read_customer
		permit get_order
	}

	delegate fraud-check-bot {
		scope "stripe/refund:amount<=500"
		ttl 24h
		ceiling inherited
	}

	ambient {
		on_exceed deny
	}

	selector account {
		source "https://api.internal/account"
		cache 30s
	}

	credential zendesk {
		backend vault
		path secret/data/zendesk
		ttl 15m
	}

	rules {
		deny! shell/*
	}
}

system global {
	version "1.0"
}

`

	p := New(src)
	doc, err := p.ParseDocument()
	if err != nil {
		t.Fatalf("ParseDocument error: %v", err)
	}
	if len(doc.Agents) != 1 {
		t.Fatalf("expected one agent, got %d", len(doc.Agents))
	}
	if doc.Agents[0].ID != "support-bot" {
		t.Fatalf("unexpected agent id: %q", doc.Agents[0].ID)
	}
	if len(doc.Agents[0].Budgets) != 1 || len(doc.Agents[0].Phases) != 1 || len(doc.Agents[0].Delegates) != 1 || len(doc.Agents[0].Selectors) != 1 || len(doc.Agents[0].Credentials) != 1 {
		t.Fatalf("structured blocks were not parsed: %+v", doc.Agents[0])
	}
	if len(doc.Systems) != 1 || doc.Systems[0].ID != "global" {
		t.Fatalf("expected top-level system block, got %+v", doc.Systems)
	}
	if len(doc.Agents[0].Rules) != 1 {
		t.Fatalf("expected rules block to survive parsing, got %d rules", len(doc.Agents[0].Rules))
	}
}
