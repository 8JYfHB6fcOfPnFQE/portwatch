# portwatch

A lightweight CLI daemon that monitors open ports and alerts on unexpected changes with configurable rules.

---

## Installation

```bash
go install github.com/youruser/portwatch@latest
```

Or build from source:

```bash
git clone https://github.com/youruser/portwatch.git && cd portwatch && go build -o portwatch .
```

---

## Usage

Start the daemon with a config file:

```bash
portwatch --config config.yaml
```

Example `config.yaml`:

```yaml
interval: 30s
alert:
  method: log
rules:
  allow:
    - 22
    - 80
    - 443
  deny: all
```

portwatch will poll open ports at the configured interval and emit an alert whenever a port outside your ruleset is detected.

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--config` | `config.yaml` | Path to config file |
| `--interval` | `30s` | Poll interval |
| `--verbose` | `false` | Enable verbose logging |

---

## How It Works

1. Reads the current list of open TCP/UDP ports at each interval
2. Compares against the configured allow/deny rules
3. Fires an alert (log, webhook, or stdout) on any unexpected change

---

## License

MIT © 2024 youruser