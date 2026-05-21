package formatter

import (
	"bytes"
	"strings"
	"testing"
)

func TestParseFormat(t *testing.T) {
	tests := []struct {
		input   string
		want    Format
		wantErr bool
	}{
		{"raw", FormatRaw, false},
		{"", FormatRaw, false},
		{"numbered", FormatNumbered, false},
		{"json", FormatJSON, false},
		{"JSON", FormatJSON, false},
		{"NUMBERED", FormatNumbered, false},
		{"csv", FormatRaw, true},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := ParseFormat(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ParseFormat(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("ParseFormat(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestWriteLineRaw(t *testing.T) {
	var buf bytes.Buffer
	fw := NewWriter(&buf, FormatRaw)
	if err := fw.WriteLine(5, "hello world"); err != nil {
		t.Fatal(err)
	}
	if got := buf.String(); got != "hello world\n" {
		t.Errorf("got %q, want %q", got, "hello world\n")
	}
}

func TestWriteLineNumbered(t *testing.T) {
	var buf bytes.Buffer
	fw := NewWriter(&buf, FormatNumbered)
	if err := fw.WriteLine(42, "some log line"); err != nil {
		t.Fatal(err)
	}
	if got := buf.String(); got != "42\tsome log line\n" {
		t.Errorf("got %q, want %q", got, "42\tsome log line\n")
	}
}

func TestWriteLineJSON(t *testing.T) {
	var buf bytes.Buffer
	fw := NewWriter(&buf, FormatJSON)
	if err := fw.WriteLine(3, `say "hello"`); err != nil {
		t.Fatal(err)
	}
	want := `{"line":3,"text":"say \"hello\""}` + "\n"
	if got := buf.String(); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestWriteLines(t *testing.T) {
	var buf bytes.Buffer
	fw := NewWriter(&buf, FormatNumbered)
	lines := []string{"alpha", "beta", "gamma"}
	if err := fw.WriteLines(lines, 10); err != nil {
		t.Fatal(err)
	}
	got := buf.String()
	for i, want := range []string{"10\talpha\n", "11\tbeta\n", "12\tgamma\n"} {
		if !strings.Contains(got, want) {
			t.Errorf("line %d: output missing %q in:\n%s", i, want, got)
		}
	}
}
