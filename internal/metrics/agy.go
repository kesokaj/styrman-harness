package metrics

import (
	"fmt"
	"io"
)

// ingestAgy reads Antigravity's conversation databases (plan Phase 3). Until
// then it says so rather than reporting an empty Antigravity as read.
func (s *Store) ingestAgy(o Options, r *Result, log io.Writer) {
	fmt.Fprintln(log, "metrics: Antigravity sessions are not read yet (plan 2026-10-08-metrics-store, Phase 3)")
}
