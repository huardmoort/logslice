package filter

import (
	"regexp"
	"strings"
)

// Filter holds criteria for matching log lines beyond time range.
type Filter struct {
	pattern   *regexp.Regexp
	keywords  []string
	exclude   *regexp.Regexp
}

// Options configures a Filter.
type Options struct {
	// Pattern is a regex that a line must match to be included.
	Pattern string
	// Keywords are plain-text substrings; a line must contain ALL of them.
	Keywords []string
	// Exclude is a regex; matching lines are dropped.
	Exclude string
}

// New creates a Filter from Options. Returns an error if any regex is invalid.
func New(opts Options) (*Filter, error) {
	f := &Filter{}

	if opts.Pattern != "" {
		re, err := regexp.Compile(opts.Pattern)
		if err != nil {
			return nil, err
		}
		f.pattern = re
	}

	if opts.Exclude != "" {
		re, err := regexp.Compile(opts.Exclude)
		if err != nil {
			return nil, err
		}
		f.exclude = re
	}

	f.keywords = opts.Keywords
	return f, nil
}

// Match reports whether line passes all filter criteria.
func (f *Filter) Match(line string) bool {
	if f.exclude != nil && f.exclude.MatchString(line) {
		return false
	}

	if f.pattern != nil && !f.pattern.MatchString(line) {
		return false
	}

	for _, kw := range f.keywords {
		if !strings.Contains(line, kw) {
			return false
		}
	}

	return true
}

// IsEmpty returns true when the filter has no criteria (matches everything).
func (f *Filter) IsEmpty() bool {
	return f.pattern == nil && f.exclude == nil && len(f.keywords) == 0
}
