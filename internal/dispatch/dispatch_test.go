package dispatch

import (
	"testing"

	"github.com/kesokaj/styrman-harness/internal/trail"
)

func TestSelectRows(t *testing.T) {
	owed := []trail.Row{trail.RowVerifier, trail.RowChangelog, trail.RowReviewer}

	t.Run("empty selects everything", func(t *testing.T) {
		got, err := selectRows(owed, "")
		if err != nil || len(got) != 3 {
			t.Fatalf("got %v, %v", got, err)
		}
	})

	t.Run("narrows and preserves dispatch order", func(t *testing.T) {
		got, err := selectRows(owed, "reviewer,verifier")
		if err != nil {
			t.Fatal(err)
		}
		// Verifier must stay first: it is cheapest and can make the rest moot.
		if len(got) != 2 || got[0] != trail.RowVerifier || got[1] != trail.RowReviewer {
			t.Errorf("got %v, want [verifier reviewer]", got)
		}
	})

	t.Run("naming an unowed row is an error", func(t *testing.T) {
		if _, err := selectRows(owed, "security"); err == nil {
			t.Error("selecting a row that is not owed was silently accepted")
		}
	})

	t.Run("tolerates whitespace and empties", func(t *testing.T) {
		got, err := selectRows(owed, " changelog , ")
		if err != nil || len(got) != 1 || got[0] != trail.RowChangelog {
			t.Errorf("got %v, %v", got, err)
		}
	})
}

func TestParseRunToleratesUnknownShapes(t *testing.T) {
	t.Run("claude envelope", func(t *testing.T) {
		out := parseRun([]byte(`{"result":"done","total_cost_usd":0.25,"is_error":false}`))
		if out.Result != "done" || out.CostUSD != 0.25 {
			t.Errorf("got %+v", out)
		}
	})
	t.Run("not json at all", func(t *testing.T) {
		out := parseRun([]byte("plain text output"))
		if out.Result != "plain text output" {
			t.Errorf("got %+v", out)
		}
	})
	t.Run("missing cost field is zero, not an error", func(t *testing.T) {
		out := parseRun([]byte(`{"response":"ok"}`))
		if out.Result != "ok" || out.CostUSD != 0 {
			t.Errorf("got %+v", out)
		}
	})
}
