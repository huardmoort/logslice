# logslice

Fast log file slicer that extracts time-range segments from large structured log files.

## Installation

```bash
go install github.com/yourusername/logslice@latest
```

## Usage

Extract log entries between two timestamps:

```bash
logslice --from "2024-01-15T08:00:00Z" --to "2024-01-15T09:00:00Z" --file app.log
```

Pipe output to a new file:

```bash
logslice --from "2024-01-15T08:00:00Z" --to "2024-01-15T09:00:00Z" --file app.log > slice.log
```

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--file` | Path to the log file | stdin |
| `--from` | Start timestamp (RFC3339) | required |
| `--to` | End timestamp (RFC3339) | required |
| `--format` | Timestamp format layout | RFC3339 |

### Example Output

```
2024-01-15T08:02:11Z INFO  server started on :8080
2024-01-15T08:14:33Z INFO  request received method=GET path=/health
2024-01-15T08:57:42Z ERROR database connection timeout retries=3
```

## How It Works

logslice uses binary search to efficiently locate the start of the requested time range, avoiding the need to scan the entire file. This makes it significantly faster than `grep` on large log files.

## License

MIT — see [LICENSE](LICENSE) for details.