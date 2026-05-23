// Package aggregate provides lightweight aggregation of field values
// extracted from log lines.
//
// Two modes are supported:
//
//   - count  – tracks how many times each distinct value appears.
//   - unique – collects the set of distinct values without counts.
//
// Typical usage:
//
//	agg, err := aggregate.New(aggregate.ModeCount)
//	if err != nil { … }
//
//	for _, line := range lines {
//	    val, _ := extractor.Extract(line)
//	    agg.Add(val)
//	}
//
//	for _, result := range agg.Results() {
//	    fmt.Println(result)
//	}
//
// Results are returned in lexicographic order of the aggregated key.
// Call Reset to reuse the same Aggregator across multiple passes.
package aggregate
