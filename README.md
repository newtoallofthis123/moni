# moni

Personal finance CLI backed by local SQLite. Track accounts, transactions, categories, recurring items, debts, and savings buckets — all from the terminal.

Single binary. No external services. `--output json` on every read command for scripting/AI agent consumption.

## Install

```bash
go install github.com/newtoallofthis123/moni@latest
```

Or build from source:

```bash
git clone https://github.com/newtoallofthis123/moni.git
cd moni
go build -o moni .
```

## Quick start

```bash
moni init                                          # create DB + seed categories
moni account add checking --type bank              # add an account
moni add income 5000 --cat salary --note "March"   # log income
moni add expense 50 --cat food --note "Groceries"  # log expense
moni balance                                       # check balances
moni transactions                                  # list transactions
moni summary                                       # monthly summary
```

## Commands

### Core

| Command | Purpose |
|---|---|
| `moni init` | Create DB + seed default categories |
| `moni account add <name> --type <type>` | Add account (bank/cash/credit/wallet/other) |
| `moni account list` | List accounts |
| `moni account edit <name> --name <new> --type <type>` | Edit account |
| `moni balance` | Show all account balances |
| `moni add expense <amount> --cat <cat> --note <desc>` | Log expense |
| `moni add income <amount> --cat <cat> --note <desc>` | Log income |
| `moni transactions [--cat <cat>] [--since <period>]` | List transactions |
| `moni transaction delete <id>` | Delete transaction (reverses balance) |
| `moni category add <name>` | Add category |
| `moni category list` | List categories |

### Debts & People

| Command | Purpose |
|---|---|
| `moni person add <name> [--phone <phone>]` | Add person |
| `moni person list` | List people |
| `moni person history <name>` | Person's transactions & debts |
| `moni debt add <person> <amount> <i_owe\|they_owe> [--note]` | Record debt |
| `moni debt settle <person> <amount>` | Settle debt (FIFO) |
| `moni debt list` | Show open debts |
| `moni debt delete <id>` | Delete debt |
| `moni link <txn_id> --persons <names...> [--note]` | Link transaction to people |

### Buckets & Recurring

| Command | Purpose |
|---|---|
| `moni bucket create <name> --target <amount>` | Create savings bucket |
| `moni bucket add <name> <amount>` | Allocate to bucket |
| `moni bucket status [<name>]` | Show bucket progress |
| `moni bucket edit <name> --name <new> --target <amount>` | Edit bucket |
| `moni bucket delete <name>` | Delete bucket |
| `moni recurring add <desc> <amount> --cat <cat> --every <freq> --due <day>` | Add recurring item |
| `moni recurring list` | List active recurring items |
| `moni recurring delete <id>` | Deactivate recurring item |
| `moni summary [--month YYYY-MM]` | Spending/income summary |

### Global flags

- `--output text|table|json` on all read commands (default: `table`)

## Development

```bash
just build    # go build -o moni .
just test     # go test ./...
just lint     # go vet ./...
just run ...  # go run . <args>
```

## Tech

- Go + [Cobra](https://github.com/spf13/cobra) CLI
- [modernc.org/sqlite](https://pkg.go.dev/modernc.org/sqlite) (pure Go, no CGO)
- Data stored in `~/.moni/moni.db`

## License

MIT
