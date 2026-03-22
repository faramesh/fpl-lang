package parser

import "testing"

func TestParseDocumentOK(t *testing.T) {
	src := `agent payment-bot {
  default deny
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
	if doc.Agents[0].Name != "payment-bot" {
		t.Fatalf("unexpected agent name: %q", doc.Agents[0].Name)
	}
	ag := doc.Agents[0]
	if ag.Default != "deny" {
		t.Fatalf("unexpected default: %q", ag.Default)
	}
	if len(ag.Rules) != 2 {
		t.Fatalf("expected two rules, got %d", len(ag.Rules))
	}
	if ag.Rules[0].Effect != "deny!" || ag.Rules[0].Tool != "shell/run" {
		t.Fatalf("unexpected first rule: %+v", ag.Rules[0])
	}
	if ag.Rules[1].Condition == nil {
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
