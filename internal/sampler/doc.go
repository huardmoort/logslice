// Package sampler implements log-line sampling strategies for logslice.
//
// Two modes are supported:
//
//	"nth"    – deterministic: keeps every Nth line (rate=N).
//	           Rate=1 keeps all lines; rate=10 keeps 1 in 10.
//
//	"random" – probabilistic: each line is independently kept
//	           with probability 1/rate.  Useful when you need an
//	           unbiased sample but do not care about regularity.
//
// Usage:
//
//	s, err := sampler.New("nth", 10, 0)
//	for scanner.Scan() {
//	    if s.Keep() {
//	        process(scanner.Text())
//	    }
//	}
package sampler
