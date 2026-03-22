# Getting Started with FPL

Write your first AI agent governance policy in five minutes.

## Prerequisites

Install Faramesh (FPL ships with the CLI):

```bash
brew install faramesh/tap/faramesh
```

Or:

```bash
curl -fsSL https://raw.githubusercontent.com/faramesh/faramesh-core/main/install.sh | bash
```

Confirm it works:

```bash
faramesh version
```

## Step 1: Write a policy

Create a file called `policy.fpl`:

```fpl
agent my-agent {
  default deny

  rules {
    deny! shell/run
      when cmd matches "rm -rf|DROP TABLE|shutdown"
      reason: "destructive command blocked"

    defer stripe/refund
      when amount > 500
      notify: "finance-team"
      reason: "high value refund needs approval"

    permit http/get
    permit read_customer
  }
}
```

This policy:
- Permanently blocks destructive shell commands (the `deny!` cannot be overridden)
- Routes high-value refunds to a human for approval
- Allows safe read operations
- Denies everything else (the `default deny` at the top)

## Step 2: Validate

Check the policy for syntax and semantic errors:

```bash
faramesh policy validate policy.fpl
```

You should see:

```
✓ policy.fpl is valid (4 rules, 1 mandatory deny)
```

## Step 3: Run your agent under governance

Prepend `faramesh run` to your normal agent command:

```bash
faramesh run --policy policy.fpl -- python agent.py
```

Faramesh will:
1. Auto-detect the agent framework (LangChain, CrewAI, etc.)
2. Patch the framework's tool dispatch to route through governance
3. Strip ambient credentials from the environment
4. Start the governance daemon
5. Run your agent

## Step 4: Watch decisions

Open a second terminal and stream live verdicts:

```bash
faramesh audit tail
```

You'll see output like:

```
[10:00:15] PERMIT  http/get           url=api.example.com     latency=8ms
[10:00:17] DENY    shell/run          cmd="rm -rf /"          policy=deny!
[10:00:18] PERMIT  read_customer      id=cust_abc             latency=9ms
[10:00:20] DEFER   stripe/refund      amount=$12,000          awaiting approval
```

## Step 5: Handle deferred actions

When a tool call is deferred, approve or deny it:

```bash
faramesh agent approve <defer-token>
faramesh agent deny <defer-token>
```

## Next steps

- [Language Reference](LANGUAGE_REFERENCE.md) — every keyword and syntax construct
- [Examples](../examples/) — ready-to-use policies for common agent types
- [Comparison](COMPARISON.md) — how FPL compares to OPA, Cedar, and YAML
