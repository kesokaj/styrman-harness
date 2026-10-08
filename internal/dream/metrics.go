package dream

import (
	"fmt"
	"io"

	"github.com/kesokaj/styrman-harness/internal/metrics"
)

// ingestMetrics reads every session's shape into the metrics store before
// the dream (R3, ADR-055): store-wide, whatever --project names, so nothing
// expires unread. It returns ingest's summary line; a source that failed is
// an error carrying that line. A variable so a test can make it fail.
var ingestMetrics = func(home string) (string, error) {
	s, err := metrics.Open(metrics.Path(home))
	if err != nil {
		return "", err
	}
	defer s.Close()
	r, err := metrics.Ingest(s, metrics.Options{Home: home}, io.Discard)
	if err != nil {
		return "", err
	}
	if len(r.Failed) > 0 {
		return "", fmt.Errorf("%s", r.Line())
	}
	return r.Line(), nil
}

// ingestFirst runs ingestMetrics and logs its outcome. It never fails the
// dream: the metrics are a side record, and a dream that stopped for them
// would also stop the proposals (D16).
func (o Options) ingestFirst(logf func(string)) {
	line, err := ingestMetrics(o.Home)
	if err != nil {
		logf("metrics ingest failed: " + scrub(err.Error(), o.Home))
		return
	}
	logf(line)
}
