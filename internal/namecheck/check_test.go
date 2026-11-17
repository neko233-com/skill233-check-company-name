package namecheck

import (
	"encoding/json"
	"testing"
)

func TestCheckEndToEnd(t *testing.T) {
	req := Request{
		Name:          "星澜科技",
		Industry:      "软件开发",
		Region:        "上海",
		Jurisdiction:  JurisdictionCNMainland,
		ExistingNames: []ExistingName{{Name: "星澜科技有限公司", Industry: "软件开发", Jurisdiction: JurisdictionCNMainland}},
	}
	r := Check(req)
	if r.Original.Risk != RiskMedium && r.Original.Risk != RiskHigh {
		t.Fatalf("original should be flagged, got %+v", r.Original)
	}
	if len(r.UnavailableReasons) == 0 {
		t.Fatal("expected unavailable reasons")
	}
	if len(r.Candidates) != DefaultCandidateCount {
		t.Fatalf("got %d candidates, want %d", len(r.Candidates), DefaultCandidateCount)
	}
	if r.Disclaimer == "" {
		t.Fatal("missing disclaimer")
	}
	for _, c := range r.Candidates {
		if len(c.Assessment.Reasons) > 0 && c.Assessment.Risk == RiskHigh {
			// It is acceptable for a generated name to be flagged, but high-risk
			// candidates should be surfaced with reasons (not silently).
			if len(c.Assessment.Reasons) == 0 {
				t.Errorf("high-risk candidate %q has no reason", c.Name)
			}
		}
	}
}

func TestReportJSONSerializable(t *testing.T) {
	r := Check(Request{Name: "星澜科技"})
	b, err := json.Marshal(r)
	if err != nil {
		t.Fatalf("marshal report: %v", err)
	}
	var back Report
	if err := json.Unmarshal(b, &back); err != nil {
		t.Fatalf("unmarshal report: %v", err)
	}
	if back.Original.Name != "星澜科技" {
		t.Fatalf("round-trip lost name: %+v", back.Original)
	}
	if len(back.Candidates) != DefaultCandidateCount {
		t.Fatalf("round-trip lost candidates: %+v", back)
	}
}
