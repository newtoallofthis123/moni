# moni MVP — Personal Finance CLI

**Ref:** MVP Research (internal)

## Goal

Ship a working CLI (`moni`) backed by a local SQLite database that lets you manually encode your financial life — accounts, transactions, categories, recurring items, debts, and savings buckets. All commands must support `--output json` for AI agent (Lily) consumption.

## Schema

Implement all 7 tables exactly as specified in the research doc:
- `accounts`, `transactions`, `categories`, `recurring`, `buckets`, `persons`, `debts`, `transaction_persons`

SQLite, single file (`~/.moni/moni.db`). Run migrations on first use via `moni init`.

## Commands

### Core (must-have)

| Command | Purpose |
|---|---|
| `moni init` | Create DB + seed default categories |
| `moni account add <name> --type <type>` | Add account |
| `moni account list` | List accounts |
| `moni balance` | Show all account balances |
| `moni add expense <amount> --cat <cat> --note <desc> [--account <acct>]` | Log outgoing transaction |
| `moni add income <amount> --cat <cat> --note <desc> [--account <acct>]` | Log incoming transaction |
| `moni transactions [--cat <cat>] [--since <period>]` | List transactions with filters |
| `moni category add <name>` | Add category |
| `moni category list` | List categories |

### Secondary (must-have for MVP)

| Command | Purpose |
|---|---|
| `moni debt add <person> <amount> <i_owe\|they_owe> [--note]` | Record a debt |
| `moni debt settle <person> <amount>` | Settle debt |
| `moni debt list` | Show open debts |
| `moni bucket create <name> --target <amount>` | Create savings bucket |
| `moni bucket add <name> <amount>` | Allocate to bucket |
| `moni bucket status [<name>]` | Show bucket progress |
| `moni recurring add <desc> <amount> --cat <cat> --every <freq> --due <day>` | Add recurring item |
| `moni recurring list` | List active recurring items |
| `moni link <txn_id> --persons <names...> [--note]` | Link transaction to people |
| `moni person add <name> [--phone]` | Add person |
| `moni person history <name>` | Show person's transactions & debts |
| `moni summary [--month <month>]` | Spending/income summary |

### Global flags

- `--output text|table|json` on all read commands
- `--account` defaults to first account if omitted on write commands

## Repo

`github.com/newtoallofthis123/moni`

## Tech

- **Language:** Go
- **CLI framework:** `cobra`
- **SQLite:** `modernc.org/sqlite` (pure Go, no CGO)
- **No external services.** Pure local. Single binary.
- **Build:** `go build -o moni .` — single binary, no runtime deps
- **Structure:**
  - `cmd/` — cobra command definitions (one file per command group)
  - `internal/db/` — connection, migrations, queries
  - `internal/models/` — structs per entity
  - `internal/format/` — text/table/json output
- Use just 
- Use proper tests

## Out of scope

Per research doc: investments/portfolio, budgets, bank import/CSV, multi-currency.

## Done when

1. `moni init` creates the DB with all tables
2. Every command listed above works end-to-end
3. `--output json` works on all read commands
4. A second user (or Lily) can install and use it with zero config beyond `moni init`
