package parser

import (
	"testing"
	"time"
)

func TestParseTimestamp(t *testing.T) {
	cases := []struct {
		input      string
		wantLayout string
		wantErr    bool
	}{
		{"2024-03-15T12:34:56Z", time.RFC3339, false},
		{"2024-03-15T12:34:56.123456789Z", time.RFC3339Nano, false},
		{"2024-03-15 12:34:56", "2006-01-02 15:04:05", false},
		{"2024-03-15 12:34:56.000", "2006-01-02 15:04:05.000", false},
		{"15/Mar/2024:12:34:56 +0000", "02/Jan/2006:15:04:05 -0700", false},
		{"not-a-timestamp", "", true},
		{"", "", true},
	}

	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			_, layout, err := ParseTimestamp(tc.input)
			if tc.wantErr {
				if err == nil {
					t.Errorf("expected error for input %q, got nil", tc.input)
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error for input %q: %v", tc.input, err)
				return
			}
			if layout != tc.wantLayout {
				t.Errorf("layout mismatch for %q: got %q, want %q", tc.input, layout, tc.wantLayout)
			}
		})
	}
}

func TestInRange(t *testing.T) {
	base := time.Date(2024, 3, 15, 12, 0, 0, 0, time.UTC)
	start := base
	end := base.Add(2 * time.Hour)

	cases := []struct {
		name string
		t    time.Time
		want bool
	}{
		{"before range", base.Add(-1 * time.Second), false},
		{"at start", start, true},
		{"in middle", base.Add(1 * time.Hour), true},
		{"at end", end, true},
		{"after range", end.Add(1 * time.Second), false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := InRange(tc.t, start, end)
			if got != tc.want {
				t.Errorf("InRange(%v) = %v, want %v", tc.t, got, tc.want)
			}
		})
	}
}
