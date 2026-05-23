// Package redact implements sensitive-data redaction for log lines.
//
// A Redactor holds a list of compiled regular expressions and a
// placeholder string. When Apply is called on a log line, every
// substring matched by any pattern is replaced with the placeholder,
// ensuring that tokens, passwords, IP addresses, or other sensitive
// values are scrubbed before the line reaches any output sink.
//
// Usage:
//
//	r, err := redact.New(
//		[]string{`password=\S+`, `token=\S+`},
//		"[REDACTED]",
//	)
//	if err != nil {
//		log.Fatal(err)
//	}
//	clean := r.Apply(rawLine)
//
// An empty pattern list produces a no-op Redactor; callers can check
// IsNoop to skip the Apply call entirely in hot paths.
package redact
