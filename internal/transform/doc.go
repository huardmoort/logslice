// Package transform provides composable line transformation primitives
// for use in the logslice processing pipeline.
//
// A Chain holds an ordered sequence of Func transformations. Each Func
// receives the current line and returns the (possibly modified) line along
// with a boolean indicating whether the line should continue through the
// pipeline. Returning false drops the line immediately; no subsequent
// functions in the chain are called.
//
// Built-in transformations:
//
//	- TrimSpace  – strip leading/trailing whitespace
//	- DropEmpty  – discard blank lines
//	- ReplaceAll – substitute substrings
//	- MaxLength  – hard-truncate long lines
//	- PrependField / AppendField – inject key=value metadata
//	- AddLineNumber – prefix lines with a sequential counter
//
// Example:
//
//	chain := transform.New(
//		transform.TrimSpace(),
//		transform.DropEmpty(),
//		transform.PrependField("source", filename),
//	)
//	for _, line := range lines {
//		if out, ok := chain.Apply(line); ok {
//			fmt.Println(out)
//		}
//	}
package transform
